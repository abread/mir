package general

import (
	"time"

	"golang.org/x/exp/slices"
)

type estimator struct {
	estimates []time.Duration
	headIdx   int
	tailIdx   int
	len       int
}

func NewEstimator(windowSize int) *estimator {
	return &estimator{
		estimates: make([]time.Duration, windowSize),
	}
}

func (e *estimator) AddSample(sample time.Duration) {
	e.estimates[e.tailIdx] = sample
	e.tailIdx = (e.tailIdx + 1) % len(e.estimates)
	if e.len < len(e.estimates) {
		e.len++
	}

	if e.len == len(e.estimates) {
		e.headIdx = (e.headIdx + 1) % len(e.estimates)
	}
}

func (e estimator) sortedSamples() []time.Duration {
	s := slices.Clone(e.estimates[:e.len])
	slices.Sort(s)
	return s
}

func (e estimator) MaxEstimate() time.Duration {
	if e.len == 0 {
		// the maximum representable duration is not a good fit because it could lead to overflows
		// return a ridiculously high, but not almost overflowing estimate
		return 24 * time.Hour
	}

	// P66 (accounts for ~1/3 nodes being byzantine and delaying operations to artificially inflate estimates)
	s := e.sortedSamples()
	idx := len(s) * 66 / 100

	return s[idx]
}

func (e estimator) MinEstimate() time.Duration {
	if e.len == 0 {
		return 0
	}

	// P10
	s := e.sortedSamples()
	idx := len(s) * 10 / 100

	return s[idx]
}
