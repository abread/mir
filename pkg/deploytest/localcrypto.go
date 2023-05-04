package deploytest

import (
	"fmt"

	mirCrypto "github.com/filecoin-project/mir/pkg/crypto"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	t "github.com/filecoin-project/mir/pkg/types"
)

type LocalCryptoSystem interface {
	Crypto(id t.NodeID) (mirCrypto.Crypto, error)
	Module(id t.NodeID) (modules.Module, error)
}

type localPseudoCryptoSystem struct {
	nodeIDs    []t.NodeID
	cryptoType string
}

// NewLocalCryptoSystem creates an instance of LocalCryptoSystem suitable for tests.
// In the current implementation, cryptoType can only be "pseudo" or "dummy".
func NewLocalCryptoSystem(cryptoType string, nodeIDs []t.NodeID, _ logging.Logger) LocalCryptoSystem {
	return &localPseudoCryptoSystem{nodeIDs, cryptoType}
}

func (cs *localPseudoCryptoSystem) Crypto(id t.NodeID) (mirCrypto.Crypto, error) {
	if cs.cryptoType == "pseudo" {
		return mirCrypto.InsecureCryptoForTestingOnly(cs.nodeIDs, id, mirCrypto.DefaultPseudoSeed)
	} else if cs.cryptoType == "dummy" {
		return &mirCrypto.DummyCrypto{
			DummySig: []byte("a valid signature"),
		}, nil
	} else {
		return nil, fmt.Errorf("unknown local crypto system type: %v (must be pseudo or dummy)", cs.cryptoType)
	}
}

func (cs *localPseudoCryptoSystem) Module(id t.NodeID) (modules.Module, error) {
	c, err := cs.Crypto(id)
	if err != nil {
		return nil, err
	}
	return mirCrypto.New(c), nil
}
