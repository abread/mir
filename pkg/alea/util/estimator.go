package util

import (
	"time"

	"golang.org/x/exp/slices"
)

type Estimator struct {
	estimates []time.Duration
	headIdx   int
	tailIdx   int
	len       int
}

func NewEstimator(windowSize int) *Estimator {
	return &Estimator{
		estimates: make([]time.Duration, windowSize),
	}
}

func (e *Estimator) AddSample(sample time.Duration) {
	e.estimates[e.tailIdx] = sample
	e.tailIdx = (e.tailIdx + 1) % len(e.estimates)
	if e.len < len(e.estimates) {
		e.len++
	}

	if e.len == len(e.estimates) {
		e.headIdx = (e.headIdx + 1) % len(e.estimates)
	}
}

func (e Estimator) sortedSamples() []time.Duration {
	s := slices.Clone(e.estimates[:e.len])
	slices.Sort(s)
	return s
}

func (e Estimator) MaxEstimate() time.Duration {
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

func (e *Estimator) Clear() {
	e.len = 0
}

func (e Estimator) Median() time.Duration {
	if e.len == 0 {
		// the maximum representable duration is not a good fit because it could lead to overflows
		// return a ridiculously high, but not almost overflowing estimate
		return 24 * time.Hour
	}

	// P50
	s := e.sortedSamples()
	idx := len(s) / 2

	return s[idx]
}

func (e Estimator) MinEstimate() time.Duration {
	if e.len == 0 {
		return 0
	}

	// P3 (accounts for ~1/3 nodes being byzantine and skipping operations to artificially deflate estimates)3
	s := e.sortedSamples()
	idx := len(s) * 33 / 100

	return s[idx]
}
