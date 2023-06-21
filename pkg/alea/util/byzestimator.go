package util

import (
	"time"

	"golang.org/x/exp/slices"
)

type ByzEstimator struct {
	children []Estimator
}

func NewByzEstimator(windowSize int, N int) ByzEstimator {
	e := ByzEstimator{
		children: make([]Estimator, N),
	}

	for i := range e.children {
		e.children[i] = NewEstimator(windowSize)
	}

	return e
}

func (e *ByzEstimator) AddSample(i int, sample time.Duration) {
	e.children[i].AddSample(sample)
}

func (e *ByzEstimator) Clear(i int) {
	e.children[i].Clear()
}

func (e *ByzEstimator) MaxEstimate() time.Duration {
	estimates := make([]time.Duration, 0, len(e.children))
	for i := 0; i < len(e.children); i++ {
		if e.children[i].len == 0 {
			continue
		}

		estimates = append(estimates, e.children[i].MaxEstimate())
	}

	if len(estimates) == 0 {
		return e.children[0].MaxEstimate() // return default max value
	}

	// P65 of P95 - skips any byz offenders (up to 1/3 of nodes) pulling the max value up
	slices.Sort(estimates)
	idx := len(estimates) * 65 / 100

	return estimates[idx]
}

func (e *ByzEstimator) Median() time.Duration {
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

	slices.Sort(estimates)
	idx := len(estimates) * 50 / 100

	return estimates[idx]
}

func (e *ByzEstimator) MinEstimate() time.Duration {
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

	// P34 of P95 - skips any byz offenders (up to 1/3 of nodes) pulling the min value down
	slices.Sort(estimates)
	idx := len(estimates) * 34 / 100

	return estimates[idx]
}
