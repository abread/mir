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
)

type Stats struct {
	lock                sync.Mutex
	reqTimestamps       map[reqKey]time.Time
	avgLatency          float64
	timestampedRequests int
	recvdRequests       int
	deliveredRequests   int

	mempoolNewBatches    uint64
	agRoundDelivers      uint64
	agRoundFalseDelivers uint64
	bcDelivers           uint64
}

type reqKey struct {
	ClientID string
	ReqNo    uint64
}

func NewStats() *Stats {
	stats := &Stats{
		reqTimestamps: make(map[reqKey]time.Time),
	}

	return stats
}

func (s *Stats) NewRequest(req *requestpb.Request) {
	s.lock.Lock()
	k := reqKey{req.ClientId, req.ReqNo}
	s.reqTimestamps[k] = time.Now()
	s.recvdRequests++
	s.lock.Unlock()
}

func (s *Stats) Delivered(req *requestpb.Request) {
	s.lock.Lock()
	s.deliveredRequests++
	k := reqKey{req.ClientId, req.ReqNo}
	if t, ok := s.reqTimestamps[k]; ok {
		delete(s.reqTimestamps, k)
		s.timestampedRequests++
		d := time.Since(t)

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

	s.avgLatency = 0
	s.timestampedRequests = 0
	s.recvdRequests = 0
	s.deliveredRequests = 0

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
	}
	_ = w.Write(record)
}
