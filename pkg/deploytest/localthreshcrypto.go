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
	ThreshCrypto(id t.NodeID) threshcrypto.ThreshCrypto
	Module(ctx context.Context, id t.NodeID) modules.Module
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

func (cs *localPseudoThreshCryptoSystem) ThreshCrypto(id t.NodeID) threshcrypto.ThreshCrypto {
	var cryptoImpl threshcrypto.ThreshCrypto
	var err error

	if cs.cryptoType == "pseudo" {
		cryptoImpl, err = threshcrypto.HerumiTBLSPseudo(cs.nodeIDs, cs.threshold, id, mirCrypto.DefaultPseudoSeed)
	} else if cs.cryptoType == "dummy" {
		cryptoImpl = &threshcrypto.DummyCrypto{
			DummySigShareSuffix: []byte("sigshare"),
			NodeID:              id,
			DummySigFull:        []byte("fullthreshsig"),
		}
	} else {
		err = fmt.Errorf("unknown local crypto system type: %v (must be pseudo or dummy)", cs.cryptoType)
	}

	if err != nil {
		panic(fmt.Sprintf("error creating crypto module: %v", err))
	}
	return cryptoImpl
}

func (cs *localPseudoThreshCryptoSystem) Module(ctx context.Context, id t.NodeID) modules.Module {
	return threshcrypto.New(ctx, threshcrypto.DefaultModuleParams(), cs.ThreshCrypto(id))
}
