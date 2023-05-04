package deploytest

import (
	"context"
	"fmt"

	mirCrypto "github.com/filecoin-project/mir/pkg/crypto"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/threshcrypto"
	t "github.com/filecoin-project/mir/pkg/types"
)

type LocalThreshCryptoSystem interface {
	ThreshCrypto(id t.NodeID) (threshcrypto.ThreshCrypto, error)
	Module(ctx context.Context, id t.NodeID) (modules.Module, error)
}

type localPseudoThreshCryptoSystem struct {
	cryptoType string
	nodeIDs    []t.NodeID
	threshold  int
}

// NewLocalCryptoSystem creates an instance of LocalCryptoSystem suitable for tests.
// In the current implementation, cryptoType can only be "pseudo".
func NewLocalThreshCryptoSystem(cryptoType string, nodeIDs []t.NodeID, threshold int, logger logging.Logger) LocalThreshCryptoSystem {
	return &localPseudoThreshCryptoSystem{cryptoType, nodeIDs, threshold}
}

func (cs *localPseudoThreshCryptoSystem) ThreshCrypto(id t.NodeID) (threshcrypto.ThreshCrypto, error) {
	if cs.cryptoType == "pseudo" {
		return threshcrypto.HerumiTBLSPseudo(cs.nodeIDs, cs.threshold, id, mirCrypto.DefaultPseudoSeed)
	} else if cs.cryptoType == "dummy" {
		return &threshcrypto.DummyCrypto{
			DummySigShareSuffix: []byte("sigshare"),
			NodeID:              id,
			DummySigFull:        []byte("fullthreshsig"),
		}, nil
	} else {
		return nil, fmt.Errorf("unknown local crypto system type: %v (must be pseudo or dummy)", cs.cryptoType)
	}
}

func (cs *localPseudoThreshCryptoSystem) Module(ctx context.Context, id t.NodeID) (modules.Module, error) {
	c, err := cs.ThreshCrypto(id)
	if err != nil {
		return nil, err
	}

	return threshcrypto.New(ctx, threshcrypto.DefaultModuleParams(), c), nil
}
