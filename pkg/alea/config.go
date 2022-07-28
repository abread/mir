package alea

import (
	"fmt"

	t "github.com/filecoin-project/mir/pkg/types"
)

type Config struct {
	Membership []t.NodeID

	MaxBatchSize t.NumRequests

	NodeBroadcastInWindowSize int

	BroadcastOutWindowSize int

	MsgBufCapacity uint

	F int
}

func CheckConfig(config *Config) error {
	if len(config.Membership) < 3*config.F+1 {
		return fmt.Errorf("Not enough nodes to tolerate %d failures", config.F)
	}

	if config.NodeBroadcastInWindowSize <= 0 {
		return fmt.Errorf("Must be able to accept at least 1 incoming broadcast per node (configured for %d)", config.NodeBroadcastInWindowSize)
	}

	if config.BroadcastOutWindowSize <= 0 {
		return fmt.Errorf("Must be able to accept at least 1 outgoing broadcast (configured for %d)", config.BroadcastOutWindowSize)
	}

	// TODO
	panic(fmt.Errorf("TODO"))
}
