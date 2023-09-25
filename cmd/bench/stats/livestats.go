// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package stats

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"sync"
	"time"

	ageventstypes "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/types"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
)

type LiveStats struct {
	lock                    sync.RWMutex
	txTimestamps            map[txKey]time.Time
	avgLatency              float64
	timestampedTransactions int
	deliveredTransactions   int

	agInputTimestamps map[uint64]time.Time
	stalledAgStart    time.Time
	currentAgRound    uint64
	avgAgStall        float64
	cumPosAgStall     float64

	avgUnanimousAgLatency float64
	innerAbbaTimeCount    int
}

type txKey struct {
	ClientID tt.ClientID
	TxNo     tt.TxNo
}

func NewLiveStats() *LiveStats {
	return &LiveStats{
		txTimestamps:      make(map[txKey]time.Time),
		agInputTimestamps: make(map[uint64]time.Time),
	}
}

func (s *LiveStats) AgDeliver(deliver *ageventstypes.Deliver) {
	s.lock.Lock()
	s.currentAgRound = deliver.Round + 1
	if t, ok := s.agInputTimestamps[deliver.Round+1]; ok {
		stall := -time.Since(t) // negative stall: input is after deliver

		// $CA_{n+1} = CA_n + {x_{n+1} - CA_n \over n + 1}$
		s.avgAgStall += (float64(stall) - s.avgAgStall) / float64(deliver.Round+1)

		delete(s.agInputTimestamps, deliver.Round+1)
	} else {
		s.stalledAgStart = time.Now()
	}
	s.lock.Unlock()
}

func (s *LiveStats) AgInput(input *ageventstypes.InputValue) {
	s.lock.Lock()
	if s.currentAgRound == input.Round {
		stall := time.Since(s.stalledAgStart)

		s.cumPosAgStall += float64(stall)
		// $CA_{n+1} = CA_n + {x_{n+1} - CA_n \over n + 1}$
		s.avgAgStall += (float64(stall) - s.avgAgStall) / float64(input.Round)
	} else {
		s.agInputTimestamps[input.Round] = time.Now()
	}
	s.lock.Unlock()
}

func (s *LiveStats) InnerAbbaTime(t *ageventstypes.InnerAbbaRoundTime) {
	s.lock.Lock()
	s.innerAbbaTimeCount++
	unanimousLatency := t.DurationNoCoin / 3
	s.avgUnanimousAgLatency += (float64(unanimousLatency) - s.avgUnanimousAgLatency) / float64(s.innerAbbaTimeCount)
	s.lock.Unlock()
}

func (s *LiveStats) Submit(tx *trantorpbtypes.Transaction) {
	s.lock.Lock()
	k := txKey{tx.ClientId, tx.TxNo}
	s.txTimestamps[k] = time.Now()
	s.lock.Unlock()
}

func (s *LiveStats) Deliver(tx *trantorpbtypes.Transaction) {
	s.lock.Lock()
	s.deliveredTransactions++
	k := txKey{tx.ClientId, tx.TxNo}
	if t, ok := s.txTimestamps[k]; ok {
		delete(s.txTimestamps, k)
		s.timestampedTransactions++
		d := time.Since(t)

		// $CA_{n+1} = CA_n + {x_{n+1} - CA_n \over n + 1}$
		s.avgLatency += (float64(d) - s.avgLatency) / float64(s.timestampedTransactions)
	}
	s.lock.Unlock()
}

func (s *LiveStats) AvgLatency() time.Duration {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return time.Duration(s.avgLatency)
}

func (s *LiveStats) DeliveredTransactions() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.deliveredTransactions
}

func (s *LiveStats) Reset() {
	s.lock.Lock()
	s.avgLatency = 0
	s.timestampedTransactions = 0
	s.deliveredTransactions = 0
	s.lock.Unlock()
}

func (s *LiveStats) WriteCSVHeader(w *csv.Writer) {
	record := []string{
		"nrDelivered",
		"tps",
		"avgLatency",
	}
	_ = w.Write(record)
}

func (s *LiveStats) WriteCSVRecord(w *csv.Writer, d time.Duration) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	tps := float64(s.deliveredTransactions) / (float64(d) / float64(time.Second))
	record := []string{
		strconv.Itoa(s.deliveredTransactions),
		strconv.Itoa(int(tps)),
		fmt.Sprintf("%.2f", time.Duration(s.avgLatency).Seconds()),
	}
	_ = w.Write(record)
}
