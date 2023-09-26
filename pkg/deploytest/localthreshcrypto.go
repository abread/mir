package deploytest

import (
	es "github.com/go-errors/errors"

	mirCrypto "github.com/filecoin-project/mir/pkg/crypto"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/threshcrypto"
	t "github.com/filecoin-project/mir/pkg/types"
)

type LocalThreshCryptoSystem interface {
	ThreshCrypto(id t.NodeID) (threshcrypto.ThreshCrypto, error)
	Module(id t.NodeID) (modules.PassiveModule, error)
}

type localPseudoThreshCryptoSystem struct {
	cryptoType string
	nodeIDs    []t.NodeID
	threshold  int
}

// NewLocalCryptoSystem creates an instance of LocalCryptoSystem suitable for tests.
// In the current implementation, cryptoType can only be "pseudo", "dummy" or "pseudo-purego".
func NewLocalThreshCryptoSystem(cryptoType string, nodeIDs []t.NodeID, threshold int) (LocalThreshCryptoSystem, error) {
	if threshold < 0 {
		return nil, es.Errorf("negative threshold: %v", threshold)
	}

	return &localPseudoThreshCryptoSystem{cryptoType, nodeIDs, threshold}, nil
}

func (cs *localPseudoThreshCryptoSystem) ThreshCrypto(id t.NodeID) (threshcrypto.ThreshCrypto, error) {
	if cs.cryptoType == "pseudo" {
		return threshcrypto.HerumiTBLSPseudo(cs.nodeIDs, cs.threshold, id, mirCrypto.DefaultPseudoSeed)
	} else if cs.cryptoType == "pseudo-purego" {
		return threshcrypto.TBLSPseudo(cs.nodeIDs, cs.threshold, id, mirCrypto.DefaultPseudoSeed)
	} else if cs.cryptoType == "dummy" {
		return &threshcrypto.DummyCrypto{
			DummySigShareSuffix: []byte("sigshare"),
			NodeID:              id,
		}, nil
	} else {
		return nil, es.Errorf("unknown local crypto system type: %v (must be pseudo or dummy)", cs.cryptoType)
	}
}

func (cs *localPseudoThreshCryptoSystem) Module(id t.NodeID) (modules.PassiveModule, error) {
	c, err := cs.ThreshCrypto(id)
	if err != nil {
		return nil, err
	}

	return threshcrypto.New(c), nil
}
