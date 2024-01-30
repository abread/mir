package stats

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

type ClientStats struct {

	// Times of submission for in-flight transactions and the lock that guards the map.
	txTimestamps   map[txKey]time.Time
	timestampsLock sync.Mutex

	// Latency histogram. Latencies are truncated (down) to the nearest step.
	// I.e., if the step is 1 millisecond, a latency of 0.99 ms will be considered as 0.
	// key: latency
	// value: number of transactions with that latency
	LatencyHist map[time.Duration]int `json:"LatencyHistogram"`

	// Throughput history.
	// key: Time since start of measurement (in multiples of SamplingPeriod)
	// value: number of transactions delivered in the time slot
	DeliveredTxs map[time.Duration]int `json:"DeliveredTxs"`

	// Time of the start of measurements.
	startTime time.Time

	// duration of the experiment.
	duration time.Duration

	// Total number of transactions delivered.
	totalDelivered int

	latencyStep    time.Duration
	SamplingPeriod time.Duration
}

func NewClientStats(
	latencyStep time.Duration,
	samplingPeriod time.Duration,
	preInitDiscardBatchCount int,
) *ClientStats {
	return &ClientStats{
		txTimestamps:   make(map[txKey]time.Time),
		LatencyHist:    map[time.Duration]int{0: 0}, // The rest of the code can assume this map is never empty.
		DeliveredTxs:   make(map[time.Duration]int),
		latencyStep:    latencyStep,
		SamplingPeriod: samplingPeriod,
	}
}

func (cs *ClientStats) ToJSON() ([]byte, error) {
	return json.Marshal(cs)
}

func (cs *ClientStats) Start() {
	cs.timestampsLock.Lock()
	cs.start()
	cs.timestampsLock.Unlock()
}

func (cs *ClientStats) start() {
	fmt.Println("restart clientstats")

	start := time.Now()
	for _, txTime := range cs.txTimestamps {
		if txTime.Before(start) {
			start = txTime
		}
	}
	cs.startTime = start
	cs.duration = cs.SamplingPeriod

	cs.LatencyHist = map[time.Duration]int{0: 0}
	cs.DeliveredTxs = make(map[time.Duration]int)
	cs.totalDelivered = 0
	cs.fillAtDuration(time.Since(start))
}

func (cs *ClientStats) Submit(tx *trantorpbtypes.Transaction) {
	txID := txKey{tx.ClientId, tx.TxNo}
	cs.timestampsLock.Lock()
	cs.txTimestamps[txID] = time.Now()
	cs.timestampsLock.Unlock()
}

func (cs *ClientStats) DeliveredBatch() {
}

func (cs *ClientStats) Deliver(tx *trantorpbtypes.Transaction) {
	cs.timestampsLock.Lock()
	defer cs.timestampsLock.Unlock()

	// Get delivery time and latency.
	t := time.Since(cs.startTime)
	txID := txKey{tx.ClientId, tx.TxNo}
	lRaw := time.Since(cs.txTimestamps[txID])
	delete(cs.txTimestamps, txID)

	// Round values to the next lower step
	t = (t / cs.SamplingPeriod) * cs.SamplingPeriod
	l := (lRaw / cs.latencyStep) * cs.latencyStep
	if l > time.Hour {
		fmt.Printf("HUGE LATENCY (raw: %s, computed: %s steps: %d steps) at time %v\n", lRaw, l, l/cs.latencyStep, t)
	}

	// Update the statistics accordingly.
	cs.LatencyHist[l]++
	cs.totalDelivered++

	// Fill periods with no delivered transactions with explicit zeroes.
	// This is redundant, as it can always be inferred from the rest of the data,
	// but makes it a bit more convenient to work with when iterating over the (sorted) items.
	cs.fillAtDuration(t)
	cs.DeliveredTxs[t]++
}

func (cs *ClientStats) AssumeDelivered(tx *trantorpbtypes.Transaction) {
	cs.timestampsLock.Lock()
	defer cs.timestampsLock.Unlock()

	txID := txKey{tx.ClientId, tx.TxNo}
	delete(cs.txTimestamps, txID)

	// Consider transaction for throughput measurement, but not for latency (latency is distorted).
	cs.totalDelivered++
}

// Fill adds padding to DeliveredTxs. In case no transactions have been delivered for some time,
// Fill adds zero values for these time slots.
// Fill should be called once more at the very end of data collection,
// especially if no transactions have been delivered at all. In such a case, DeliveredTxs would otherwise stay empty
// and not represent the true result of data collection (namely zeroes for all the duration.)
func (cs *ClientStats) Fill() {
	cs.timestampsLock.Lock()
	cs.fillAtDuration(time.Since(cs.startTime))
	cs.timestampsLock.Unlock()
}

func (cs *ClientStats) fillAtDuration(t time.Duration) {
	for t >= cs.duration+cs.SamplingPeriod {
		cs.duration += cs.SamplingPeriod
		cs.DeliveredTxs[cs.duration] = 0
	}
}

func (cs *ClientStats) AvgThroughput() float64 {
	if cs.duration == 0 {
		return 0
	}
	return float64(cs.totalDelivered) / cs.duration.Seconds()
}

func (cs *ClientStats) AvgLatency() time.Duration {
	if cs.totalDelivered == 0 {
		return 0
	}

	totalLatency := time.Duration(0)
	for latency, numTx := range cs.LatencyHist {
		totalLatency += time.Duration(numTx) * latency
	}

	return totalLatency / time.Duration(cs.totalDelivered)
}

func (cs *ClientStats) LatencyPctile(pctile float32) time.Duration {
	txCount := 0
	var result time.Duration
	maputil.IterateSorted(cs.LatencyHist, func(latency time.Duration, numTxs int) (cont bool) {
		result = latency
		txCount += numTxs
		return float32(txCount)/float32(cs.totalDelivered) < pctile
	})
	return result
}

func (cs *ClientStats) WriteCSVHeader(w *csv.Writer) error {
	record := []string{
		"duration",
		"throughput",
		"latency_avg",
		"latency_median",
		"latency_95p",
		"latency_max",
	}
	return w.Write(record)
}

func (cs *ClientStats) WriteCSVRecord(w *csv.Writer, _ time.Duration) error {
	cs.timestampsLock.Lock()
	record := []string{
		fmt.Sprintf("%.6f", cs.duration.Seconds()),
		fmt.Sprintf("%.3f", cs.AvgThroughput()),
		fmt.Sprintf("%.6f", cs.AvgLatency().Seconds()),
		fmt.Sprintf("%.6f", cs.LatencyPctile(0.5).Seconds()),
		fmt.Sprintf("%.6f", cs.LatencyPctile(0.95).Seconds()),
		fmt.Sprintf("%.6f", cs.LatencyPctile(1.0).Seconds()),
	}
	cs.timestampsLock.Unlock()
	return w.Write(record)
}
