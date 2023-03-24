// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package stats

import (
	"encoding/csv"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	"github.com/filecoin-project/mir/pkg/pb/threshcryptopb"
)

type Stats struct {
	lock                sync.Mutex
	reqTimestamps       map[reqKey]int64
	avgLatency          float64
	timestampedRequests int
	recvdRequests       int
	deliveredRequests   int

	mempoolNewBatches    uint64
	agRoundDelivers      uint64
	agRoundFalseDelivers uint64
	bcDelivers           uint64
	threshQueueSize      int
}

type reqKey struct {
	ClientID string
	ReqNo    uint64
}

func NewStats() *Stats {
	stats := &Stats{
		reqTimestamps: make(map[reqKey]int64),
	}

	return stats
}

func (s *Stats) NewRequest(req *requestpb.Request, ts int64) {
	s.lock.Lock()
	k := reqKey{req.ClientId, req.ReqNo}
	s.reqTimestamps[k] = ts
	s.recvdRequests++
	s.lock.Unlock()
}

func (s *Stats) Delivered(req *requestpb.Request, deliverTs int64) { // nolint:stylecheck
	s.lock.Lock()
	s.deliveredRequests++
	k := reqKey{req.ClientId, req.ReqNo}
	if startTs, ok := s.reqTimestamps[k]; ok { // nolint:stylecheck
		delete(s.reqTimestamps, k)
		s.timestampedRequests++
		d := deliverTs - startTs

		// $CA_{n+1} = CA_n + {x_{n+1} - CA_n \over n + 1}$
		s.avgLatency += (float64(d) - s.avgLatency) / float64(s.timestampedRequests)
	}
	s.lock.Unlock()
}

func (s *Stats) MempoolNewBatch() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.mempoolNewBatches++
}

func (s *Stats) DeliveredAgRound(ev *agevents.Deliver) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.agRoundDelivers++

	if !ev.Decision {
		s.agRoundFalseDelivers++
	}
}

func (s *Stats) DeliveredBcSlot() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.bcDelivers++
}

func (s *Stats) ThreshCryptoEvent(ev *threshcryptopb.Event) {
	var delta int

	switch ev.Type.(type) {
	case *threshcryptopb.Event_SignShare:
	case *threshcryptopb.Event_VerifyShare:
	case *threshcryptopb.Event_VerifyFull:
	case *threshcryptopb.Event_Recover:
		delta = 1
	case *threshcryptopb.Event_SignShareResult:
	case *threshcryptopb.Event_VerifyShareResult:
	case *threshcryptopb.Event_VerifyFullResult:
	case *threshcryptopb.Event_RecoverResult:
		delta = -1
	default:
		panic("unknown threshcrypto event")
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	s.threshQueueSize += delta
}

func (s *Stats) WriteCSVHeader(w *csv.Writer) {
	record := []string{
		"ts",
		"nrReceived",
		"loadtps",
		"nrDelivered",
		"tps",
		"avgLatency",
		"memSys",
		"memStackInUse",
		"memHeapAlloc",
		"memTotalAlloc",
		"memMallocs",
		"memFrees",
		"memPauseTotalNs",
		"memNumGC",
		"mempoolNewBatches",
		"agRoundDelivers",
		"agRoundFalseDelivers",
		"bcDelivers",
		"threshQueueSize",
	}
	_ = w.Write(record)
}

func (s *Stats) WriteCSVRecordAndReset(w *csv.Writer, d time.Duration) {
	var memStats runtime.MemStats

	s.lock.Lock()

	now := time.Now().UnixMilli()
	runtime.ReadMemStats(&memStats)

	deliveredReqs := s.deliveredRequests
	recvdReqs := s.recvdRequests
	avgLatency := s.avgLatency
	newBatchCount := s.mempoolNewBatches
	agRoundDelivers := s.agRoundDelivers
	agRoundFalseDelivers := s.agRoundFalseDelivers
	bcDelivers := s.bcDelivers
	threshQueueSize := s.threshQueueSize

	s.avgLatency = 0
	s.timestampedRequests = 0
	s.recvdRequests = 0
	s.deliveredRequests = 0

	s.agRoundDelivers = 0
	s.agRoundFalseDelivers = 0
	s.bcDelivers = 0

	s.lock.Unlock()

	loadTps := float64(recvdReqs) / (float64(d) / float64(time.Second))
	tps := float64(deliveredReqs) / (float64(d) / float64(time.Second))
	record := []string{
		strconv.FormatInt(now, 10),
		strconv.Itoa(recvdReqs),
		fmt.Sprintf("%.3f", loadTps),
		strconv.Itoa(deliveredReqs),
		fmt.Sprintf("%.3f", tps),
		fmt.Sprintf("%.3f", time.Duration(avgLatency).Seconds()),
		strconv.FormatUint(memStats.Sys, 10),
		strconv.FormatUint(memStats.StackInuse, 10),
		strconv.FormatUint(memStats.HeapAlloc, 10),
		strconv.FormatUint(memStats.TotalAlloc, 10),
		strconv.FormatUint(memStats.Mallocs, 10),
		strconv.FormatUint(memStats.Frees, 10),
		strconv.FormatUint(memStats.PauseTotalNs, 10),
		strconv.FormatUint(uint64(memStats.NumGC), 10),
		strconv.FormatUint(newBatchCount, 10),
		strconv.FormatUint(agRoundDelivers, 10),
		strconv.FormatUint(agRoundFalseDelivers, 10),
		strconv.FormatUint(bcDelivers, 10),
		strconv.FormatInt(int64(threshQueueSize), 10),
	}
	_ = w.Write(record)
}
