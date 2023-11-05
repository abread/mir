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

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	abbapbtypes "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	ageventstypes "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/types"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
)

type LiveStats struct {
	lock                    sync.RWMutex
	txTimestamps            map[txKey]time.Time
	avgLatency              time.Duration
	timestampedTransactions int
	deliveredTransactions   int

	ownQueueIdx       aleatypes.QueueIdx
	bcDelivers        int
	ownBcStartedCount int
	bcStallStart      time.Time
	avgBcStall        time.Duration

	agInputTimestamps map[uint64]time.Time
	stalledAgStart    time.Time
	currentAgRound    uint64
	avgAgStall        time.Duration
	cumPosAgStall     time.Duration
	nonInstantAgCount int
	trueAgDelivers    int
	falseAgDelivers   int

	estUnanimousAgTime time.Duration
	innerAbbaTimeCount int

	avgBatchSize256  int
	formedBatchCount int
}

type txKey struct {
	ClientID tt.ClientID
	TxNo     tt.TxNo
}

func NewLiveStats(ownQueueIdx aleatypes.QueueIdx) *LiveStats {
	return &LiveStats{
		txTimestamps:      make(map[txKey]time.Time),
		agInputTimestamps: make(map[uint64]time.Time),

		ownQueueIdx: ownQueueIdx,
	}
}

func (s *LiveStats) Fill() {}

func (s *LiveStats) RequestCert() {
	var zeroTime time.Time
	s.lock.Lock()
	s.ownBcStartedCount++
	if s.bcStallStart != zeroTime {
		stall := time.Since(s.bcStallStart)

		// $CA_{n+1} = CA_n + {x_{n+1} - CA_n \over n + 1}$
		s.avgBcStall += (stall - s.avgBcStall) / time.Duration(s.ownBcStartedCount)
	}
	s.lock.Unlock()
}

func (s *LiveStats) BcDeliver(cert *bcpbtypes.DeliverCert) {
	s.lock.Lock()
	s.bcDelivers++

	if cert.Cert.Slot.QueueIdx == s.ownQueueIdx {
		s.bcStallStart = time.Now()
	}
	s.lock.Unlock()
}

func (s *LiveStats) AgDeliver(deliver *ageventstypes.Deliver) {
	s.lock.Lock()
	s.currentAgRound = deliver.Round + 1
	if t, ok := s.agInputTimestamps[deliver.Round+1]; ok {
		stall := -time.Since(t) // negative stall: input is before deliver of the previous round

		// $CA_{n+1} = CA_n + {x_{n+1} - CA_n \over n + 1}$
		s.avgAgStall += (stall - s.avgAgStall) / time.Duration(s.trueAgDelivers+s.falseAgDelivers+1)

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

		s.cumPosAgStall += stall
		// $CA_{n+1} = CA_n + {x_{n+1} - CA_n \over n + 1}$
		count := s.trueAgDelivers + s.falseAgDelivers
		if count == 0 {
			count = 1
		}
		s.avgAgStall += (stall - s.avgAgStall) / time.Duration(count)
	} else {
		s.agInputTimestamps[input.Round] = time.Now()
	}
	s.lock.Unlock()
}

func (s *LiveStats) AbbaRoundContinue(roundNum string, _ *abbapbtypes.RoundContinue) {
	if roundNum != "0" {
		return
	}

	s.lock.Lock()
	s.nonInstantAgCount++
	s.lock.Unlock()
}

func (s *LiveStats) InnerAbbaTime(t *ageventstypes.InnerAbbaRoundTime) {
	s.lock.Lock()
	s.innerAbbaTimeCount++
	unanimousLatency := t.DurationNoCoin / 3
	s.estUnanimousAgTime += (unanimousLatency - s.estUnanimousAgTime) / time.Duration(s.innerAbbaTimeCount)
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
		s.avgLatency += (d - s.avgLatency) / time.Duration(s.timestampedTransactions)
	}
	s.lock.Unlock()
}

func (s *LiveStats) CutBatch(batchSz int) {
	s.lock.Lock()

	s.formedBatchCount++
	// batch size is multiplied by 256 to retain 8 bits of precision
	s.avgBatchSize256 += (batchSz*256 - s.avgBatchSize256) / s.formedBatchCount

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

	return s.avgLatency
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
		"avgBcStall",
		"avgBatchSize",
		"trueAgDelivers",
		"falseAgDelivers",
		"nonInstantAgCount",
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
	avgBcStall := s.avgBcStall
	avgBatchSize256 := s.avgBatchSize256
	trueAgDelivers := s.trueAgDelivers
	falseAgDelivers := s.falseAgDelivers
	nonInstantAgCount := s.nonInstantAgCount
	avgAgStall := s.avgAgStall
	cumPosAgStall := s.cumPosAgStall
	estUnanimousAgTime := s.estUnanimousAgTime

	s.avgLatency = 0
	s.timestampedTransactions = 0
	s.deliveredTransactions = 0
	s.bcDelivers = 0
	s.avgBcStall = 0
	s.ownBcStartedCount = 0
	s.avgBatchSize256 = 0
	s.formedBatchCount = 0
	s.trueAgDelivers = 0
	s.falseAgDelivers = 0
	s.nonInstantAgCount = 0
	s.avgAgStall = 0
	s.cumPosAgStall = 0
	s.estUnanimousAgTime = 0
	s.innerAbbaTimeCount = 0
	s.lock.Unlock()

	tps := float64(deliveredTxs) / (float64(d) / float64(time.Second))
	record := []string{
		fmt.Sprintf("%.3f", float64(time.Now().UnixMilli())/1000.0),
		strconv.Itoa(deliveredTxs),
		fmt.Sprintf("%.1f", tps),
		fmt.Sprintf("%.6f", avgLatency.Seconds()),
		strconv.Itoa(bcDelivers),
		fmt.Sprintf("%.6f", avgBcStall.Seconds()),
		fmt.Sprintf("%.3f", float64(avgBatchSize256)/256.0),
		strconv.Itoa(trueAgDelivers),
		strconv.Itoa(falseAgDelivers),
		strconv.Itoa(nonInstantAgCount),
		fmt.Sprintf("%.6f", avgAgStall.Seconds()),
		fmt.Sprintf("%.6f", cumPosAgStall.Seconds()),
		fmt.Sprintf("%.6f", estUnanimousAgTime.Seconds()),
	}
	return w.Write(record)
}
