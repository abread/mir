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

func NewEstimator(windowSize int) Estimator {
	return Estimator{
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

func (e *Estimator) sortedSamples() []time.Duration {
	s := slices.Clone(e.estimates[:e.len])
	slices.Sort(s)
	return s
}

func (e *Estimator) MaxEstimate() time.Duration {
	if e.len == 0 {
		return 0
	}

	// P85
	s := e.sortedSamples()
	idx := len(s) * 85 / 100

	return s[idx]
}

func (e *Estimator) Clear() {
	e.len = 0
}

func (e *Estimator) Median() time.Duration {
	if e.len == 0 {
		return 0
	}

	// P50
	s := e.sortedSamples()
	idx := len(s) / 2

	return s[idx]
}

func (e *Estimator) MinEstimate() time.Duration {
	if e.len == 0 {
		return 0
	}

	// P15
	s := e.sortedSamples()
	idx := len(s) * 15 / 100

	return s[idx]
}
