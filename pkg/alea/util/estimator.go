package util

import (
	"time"

	"golang.org/x/exp/slices"
)

type Estimator struct {
	samples       []time.Duration
	sortedSamples []time.Duration
	headIdx       int
	tailIdx       int
	len           int
}

func NewEstimator(windowSize int) Estimator {
	return Estimator{
		samples:       make([]time.Duration, windowSize),
		sortedSamples: make([]time.Duration, 0, windowSize),
	}
}

func (e *Estimator) AddSample(sample time.Duration) {
	e.samples[e.tailIdx] = sample
	e.tailIdx = (e.tailIdx + 1) % len(e.samples)
	if e.len < len(e.samples) {
		e.len++
	}

	if e.len == len(e.samples) {
		// replace value in sorted samples with new one
		toRemove := e.samples[e.headIdx]
		samplesView := e.sortedSamples
		for samplesView[0] != toRemove {
			middleIdx := len(samplesView) / 2
			if samplesView[middleIdx] > toRemove {
				samplesView = samplesView[:middleIdx+1]
			} else {
				samplesView = samplesView[middleIdx:]
			}
		}
		samplesView[0] = sample

		// update main samples ring
		e.headIdx = (e.headIdx + 1) % len(e.samples)
	} else {
		// max len not reached yet, append new sample
		e.sortedSamples = append(e.sortedSamples, sample)
	}

	// re-sort samples
	slices.Sort(e.sortedSamples)
}

func (e *Estimator) MaxEstimate() time.Duration {
	if e.len == 0 {
		return 0
	}

	// P95
	idx := len(e.sortedSamples) * 95 / 100

	return e.sortedSamples[idx]
}

func (e *Estimator) Clear() {
	e.len = 0
}

func (e *Estimator) Median() time.Duration {
	if e.len == 0 {
		return 0
	}

	// P50
	idx := len(e.sortedSamples) / 2

	return e.sortedSamples[idx]
}

func (e *Estimator) MinEstimate() time.Duration {
	if e.len == 0 {
		return 0
	}

	// P05
	idx := len(e.sortedSamples) * 5 / 100

	return e.sortedSamples[idx]
}
