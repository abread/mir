// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package stats

import (
	"encoding/csv"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"sync"
	"time"

	ageventstypes "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/types"
	threshcryptopbtypes "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/types"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"

	tt "github.com/filecoin-project/mir/pkg/trantor/types"
)

type Stats struct {
	lock                    sync.Mutex
	txTimestamps            map[txKey]int64
	avgLatency              float64
	timestampedTransactions int
	recvdTxs                int
	deliveredTxs            int

	mempoolNewBatches    uint64
	agRoundDelivers      uint64
	agRoundFalseDelivers uint64
	bcDelivers           uint64
	threshQueueSize      int
}

type txKey struct {
	ClientID tt.ClientID
	TxNo     tt.TxNo
}

func NewStats() *Stats {
	stats := &Stats{
		txTimestamps: make(map[txKey]int64),
	}

	return stats
}

func (s *Stats) NewTX(tx *trantorpbtypes.Transaction, ts int64) {
	s.lock.Lock()
	k := txKey{tx.ClientId, tx.TxNo}
	s.txTimestamps[k] = ts
	s.recvdTxs++
	s.lock.Unlock()
}

func (s *Stats) Delivered(tx *trantorpbtypes.Transaction, deliverTS int64) {
	s.lock.Lock()
	s.deliveredTxs++
	k := txKey{tx.ClientId, tx.TxNo}
	if startTS, ok := s.txTimestamps[k]; ok {
		delete(s.txTimestamps, k)
		s.timestampedTransactions++
		d := deliverTS - startTS

		// $CA_{n+1} = CA_n + {x_{n+1} - CA_n \over n + 1}$
		s.avgLatency += (float64(d) - s.avgLatency) / float64(s.timestampedTransactions)
	}
	s.lock.Unlock()
}

func (s *Stats) MempoolNewBatch() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.mempoolNewBatches++
}

func (s *Stats) DeliveredAgRound(ev *ageventstypes.Deliver) {
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

func (s *Stats) ThreshCryptoEvent(ev *threshcryptopbtypes.Event) {
	var delta int

	switch ev.Type.(type) {
	case *threshcryptopbtypes.Event_SignShare:
	case *threshcryptopbtypes.Event_VerifyShare:
	case *threshcryptopbtypes.Event_VerifyFull:
	case *threshcryptopbtypes.Event_Recover:
		delta = 1
	case *threshcryptopbtypes.Event_SignShareResult:
	case *threshcryptopbtypes.Event_VerifyShareResult:
	case *threshcryptopbtypes.Event_VerifyFullResult:
	case *threshcryptopbtypes.Event_RecoverResult:
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

	deliveredReqs := s.deliveredTxs
	recvdReqs := s.recvdTxs
	avgLatency := s.avgLatency
	newBatchCount := s.mempoolNewBatches
	agRoundDelivers := s.agRoundDelivers
	agRoundFalseDelivers := s.agRoundFalseDelivers
	bcDelivers := s.bcDelivers
	threshQueueSize := s.threshQueueSize

	if s.timestampedTransactions == 0 {
		avgLatency = math.NaN()
	}

	s.avgLatency = 0
	s.timestampedTransactions = 0
	s.recvdTxs = 0
	s.deliveredTxs = 0

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
