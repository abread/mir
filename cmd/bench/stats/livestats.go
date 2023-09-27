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
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
)

type LiveStats struct {
	lock                    sync.RWMutex
	txTimestamps            map[txKey]time.Time
	avgLatency              float64
	timestampedTransactions int
	deliveredTransactions   int

	bcDelivers int

	agInputTimestamps map[uint64]time.Time
	stalledAgStart    time.Time
	currentAgRound    uint64
	avgAgStall        float64
	cumPosAgStall     float64
	trueAgDelivers    int
	falseAgDelivers   int

	estUnanimousAgTime float64
	innerAbbaTimeCount int
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

func (s *LiveStats) Fill() {}

func (s *LiveStats) BcDeliver(_ *bcpbtypes.DeliverCert) {
	s.lock.Lock()
	s.bcDelivers++
	s.lock.Unlock()
}

func (s *LiveStats) AgDeliver(deliver *ageventstypes.Deliver) {
	s.lock.Lock()
	s.currentAgRound = deliver.Round + 1
	if t, ok := s.agInputTimestamps[deliver.Round+1]; ok {
		stall := -time.Since(t) // negative stall: input is before deliver of the previous round

		// $CA_{n+1} = CA_n + {x_{n+1} - CA_n \over n + 1}$
		s.avgAgStall += (stall.Seconds() - s.avgAgStall) / float64(deliver.Round+2)

		delete(s.agInputTimestamps, deliver.Round+1)
	} else {
		s.stalledAgStart = time.Now()
	}

	if deliver.Decision {
		s.trueAgDelivers++
	} else {
		s.falseAgDelivers++
	}
	s.lock.Unlock()
}

func (s *LiveStats) AgInput(input *ageventstypes.InputValue) {
	if input.Round == 0 {
		return // ignore first round
	}

	s.lock.Lock()
	if s.currentAgRound == input.Round {
		stall := time.Since(s.stalledAgStart)

		s.cumPosAgStall += stall.Seconds()
		// $CA_{n+1} = CA_n + {x_{n+1} - CA_n \over n + 1}$
		s.avgAgStall += (stall.Seconds() - s.avgAgStall) / float64(input.Round+1)
	} else {
		s.agInputTimestamps[input.Round] = time.Now()
	}
	s.lock.Unlock()
}

func (s *LiveStats) InnerAbbaTime(t *ageventstypes.InnerAbbaRoundTime) {
	s.lock.Lock()
	s.innerAbbaTimeCount++
	unanimousLatency := t.DurationNoCoin / 3
	s.estUnanimousAgTime += (unanimousLatency.Seconds() - s.estUnanimousAgTime) / float64(s.innerAbbaTimeCount)
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

func (s *LiveStats) AssumeDelivered(tx *trantorpbtypes.Transaction) {
	s.lock.Lock()

	// Consider transaction for throughput measurement, but not for latency (latency is distorted).
	s.deliveredTransactions++

	k := txKey{tx.ClientId, tx.TxNo}
	delete(s.txTimestamps, k)

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

func (s *LiveStats) WriteCSVHeader(w *csv.Writer) error {
	record := []string{
		"time",
		"nrDelivered",
		"tps",
		"avgLatency",
		"bcDelivers",
		"trueAgDelivers",
		"falseAgDelivers",
		"avgAgStall",
		"cumPosAgStall",
		"estUnanimousAgTime",
	}
	return w.Write(record)
}

func (s *LiveStats) WriteCSVRecord(w *csv.Writer, d time.Duration) error {
	s.lock.Lock()
	deliveredTxs := s.deliveredTransactions
	avgLatency := s.avgLatency
	bcDelivers := s.bcDelivers
	trueAgDelivers := s.trueAgDelivers
	falseAgDelivers := s.falseAgDelivers
	avgAgStall := s.avgAgStall
	cumPosAgStall := s.cumPosAgStall
	estUnanimousAgTime := s.estUnanimousAgTime

	s.avgLatency = 0
	s.timestampedTransactions = 0
	s.deliveredTransactions = 0
	s.bcDelivers = 0
	s.trueAgDelivers = 0
	s.falseAgDelivers = 0
	s.lock.Unlock()

	tps := float64(deliveredTxs) / (float64(d) / float64(time.Second))
	record := []string{
		fmt.Sprintf("%.3f", float64(time.Now().UnixMilli())/1000.0),
		strconv.Itoa(deliveredTxs),
		fmt.Sprintf("%.1f", tps),
		fmt.Sprintf("%.5f", time.Duration(avgLatency).Seconds()),
		strconv.Itoa(bcDelivers),
		strconv.Itoa(trueAgDelivers),
		strconv.Itoa(falseAgDelivers),
		fmt.Sprintf("%.5f", avgAgStall),
		fmt.Sprintf("%.5f", cumPosAgStall),
		fmt.Sprintf("%.5f", estUnanimousAgTime),
	}
	return w.Write(record)
}
