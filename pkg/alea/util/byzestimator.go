package util

import (
	"time"

	"golang.org/x/exp/slices"
)

type ByzEstimator struct {
	children []*Estimator

	maxEst time.Duration
	minEst time.Duration
	median time.Duration
}

func NewByzEstimator(windowSize int, N int) *ByzEstimator {
	e := &ByzEstimator{
		children: make([]*Estimator, N),
	}

	for i := range e.children {
		e.children[i] = NewEstimator(windowSize)
	}

	return e
}

func (e *ByzEstimator) AddSample(i int, sample time.Duration) {
	e.children[i].AddSample(sample)
	e.maxEst = -1
	e.minEst = -1
	e.median = -1
}

func (e *ByzEstimator) Clear(i int) {
	e.children[i].Clear()
	e.maxEst = -1
	e.minEst = -1
	e.median = -1
}

func (e *ByzEstimator) MaxEstimate() time.Duration {
	if e.maxEst == -1 {
		e.maxEst = e._maxEstimate()
	}

	return e.maxEst
}

func (e *ByzEstimator) MinEstimate() time.Duration {
	if e.minEst == -1 {
		e.minEst = e._minEstimate()
	}

	return e.minEst
}

func (e *ByzEstimator) Median() time.Duration {
	if e.median == -1 {
		e.median = e._median()
	}

	return e.median
}

func (e *ByzEstimator) _maxEstimate() time.Duration {
	estimates := make([]time.Duration, 0, len(e.children))
	for i := 0; i < len(e.children); i++ {
		if e.children[i].len == 0 {
			continue
		}

		estimates = append(estimates, e.children[i].MaxEstimate())
	}

	if len(estimates) == 0 {
		return e.children[0].MaxEstimate() // return default max value
	} else if len(estimates) == 1 {
		return estimates[0]
	}

	// ~P66 of P95 - skips any byz offenders (up to F of 3F+1 nodes) pulling the max value up
	slices.Sort(estimates)
	F := (len(estimates) - 1) / 3 // = F (when all nodes have estimates)
	idx := len(estimates) - F - 2

	return estimates[idx]
}

func (e *ByzEstimator) _median() time.Duration {
	estimates := make([]time.Duration, 0, len(e.children))
	for i := 0; i < len(e.children); i++ {
		if e.children[i].len == 0 {
			continue
		}

		estimates = append(estimates, e.children[i].Median())
	}

	if len(estimates) == 0 {
		return e.children[0].Median() // return default median value
	}

	// P50 of P50 - skips any byz offenders (up to F of 3F+1 nodes) pulling the median up/down
	slices.Sort(estimates)
	idx := len(estimates) / 2

	return estimates[idx]
}

func (e *ByzEstimator) _minEstimate() time.Duration {
	estimates := make([]time.Duration, 0, len(e.children))
	for i := 0; i < len(e.children); i++ {
		if e.children[i].len == 0 {
			continue
		}

		estimates = append(estimates, e.children[i].MinEstimate())
	}

	if len(estimates) == 0 {
		return e.children[0].MinEstimate() // return default min value
	}

	// ~P33 of P95 - skips any byz offenders (up to F of 3F+1 nodes) pulling the min value down
	slices.Sort(estimates)
	idx := (len(estimates) - 1) / 3 // = F

	return estimates[idx]
}
