package alea

import (
	"fmt"

	t "github.com/hyperledger-labs/mirbft/pkg/types"
)

type Config struct {
	Membership []t.NodeID

	MaxBatchSize t.NumRequests

	MsgBufCapacity uint
}

func CheckConfig(config *Config) error {
	// TODO
	panic(fmt.Errorf("TODO"))
}
