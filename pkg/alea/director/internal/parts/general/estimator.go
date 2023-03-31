package general

import "time"

type estimator time.Duration

func (e *estimator) AddSample(sample time.Duration) {
	firstEst := *e == 0
	*e += estimator(sample)
	if firstEst {
		*e /= 2
	}
}

func (e estimator) Estimate() time.Duration {
	return time.Duration(e)
}
