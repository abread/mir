package deploytest

import (
	"fmt"

	mirCrypto "github.com/filecoin-project/mir/pkg/crypto"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	t "github.com/filecoin-project/mir/pkg/types"
)

type LocalCryptoSystem interface {
	Crypto(id t.NodeID) mirCrypto.Crypto
	Module(id t.NodeID) modules.Module
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

func (cs *localPseudoCryptoSystem) Crypto(id t.NodeID) mirCrypto.Crypto {
	var cryptoImpl mirCrypto.Crypto
	var err error

	if cs.cryptoType == "pseudo" {
		cryptoImpl, err = mirCrypto.InsecureCryptoForTestingOnly(cs.nodeIDs, id, mirCrypto.DefaultPseudoSeed)
	} else if cs.cryptoType == "dummy" {
		cryptoImpl = &mirCrypto.DummyCrypto{
			DummySig: []byte("a valid signature"),
		}
	} else {
		err = fmt.Errorf("unknown local crypto system type: %v (must be pseudo or dummy)", cs.cryptoType)
	}

	if err != nil {
		panic(fmt.Sprintf("error creating crypto module: %v", err))
	}
	return cryptoImpl
}

func (cs *localPseudoCryptoSystem) Module(id t.NodeID) modules.Module {
	return mirCrypto.New(cs.Crypto(id))
}
