package util

import (
	"time"

	"golang.org/x/exp/slices"
)

type Estimator struct {
	samples       []time.Duration
	sortedSamples []time.Duration
	tailIdx       int
	len           int
}

func NewEstimator(windowSize int) *Estimator {
	return &Estimator{
		samples:       make([]time.Duration, windowSize),
		sortedSamples: make([]time.Duration, 0, windowSize),
	}
}

func (e *Estimator) AddSample(sample time.Duration) {
	removedValue := e.samples[e.tailIdx]
	e.samples[e.tailIdx] = sample
	e.tailIdx = (e.tailIdx + 1) % len(e.samples)

	if e.len == len(e.samples) {
		// replace value in sorted samples with new one
		samplesView := e.sortedSamples
		for samplesView[0] != removedValue {
			if len(samplesView) == 1 {
				panic("sample not found in sorted samples")
			}

			middleIdx := len(samplesView) / 2
			if samplesView[middleIdx] > removedValue {
				samplesView = samplesView[:middleIdx+1]
			} else {
				samplesView = samplesView[middleIdx:]
			}
		}
		samplesView[0] = sample
	} else {
		// max len not reached yet, append new sample
		e.sortedSamples = append(e.sortedSamples, sample)
	}

	if e.len < len(e.samples) {
		e.len++
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
