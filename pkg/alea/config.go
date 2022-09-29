package alea

import (
	"fmt"

	t "github.com/filecoin-project/mir/pkg/types"
)

type Config struct {
	Membership []t.NodeID

	MaxBatchSize t.NumRequests

	NodeBroadcastWindowSize int

	AgreementWindowSize int

	F int

	NetModuleName          t.ModuleID
	ThreshCryptoModuleName t.ModuleID
	AleaModuleName         t.ModuleID
	VCBCModuleFactoryName  t.ModuleID
	ABAModuleFactoryName   t.ModuleID
}

func CheckConfig(config *Config) error {
	if len(config.Membership) < 3*config.F+1 {
		return fmt.Errorf("Not enough nodes to tolerate %d faulty ones (with 3f+1 nodes, up to f can be byzantine)", config.F)
	}

	if config.NodeBroadcastWindowSize <= 0 {
		return fmt.Errorf("Must be able to accept at least 1 broadcast per node (configured for %d)", config.NodeBroadcastWindowSize)
	}

	if config.AgreementWindowSize <= 0 {
		return fmt.Errorf("Must be able to accept at least 1 ABBA instance (configured for %d)", config.AgreementWindowSize)
	}

	// TODO
	panic(fmt.Errorf("TODO"))
}
