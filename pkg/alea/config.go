package alea

import (
	"fmt"

	t "github.com/filecoin-project/mir/pkg/types"
)

type Config struct {
	Membership []t.NodeID

	MaxBatchSize t.NumRequests

	MsgBufCapacity uint

	F int
}

func CheckConfig(config *Config) error {
	if len(config.Membership) < 3*config.F+1 {
		return fmt.Errorf("Not enough nodes to tolerate %d failures", config.F)
	}

	// TODO
	panic(fmt.Errorf("TODO"))
}
