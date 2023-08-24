package deploytest

import (
	es "github.com/go-errors/errors"

	mirCrypto "github.com/filecoin-project/mir/pkg/crypto"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	t "github.com/filecoin-project/mir/pkg/types"
)

type LocalCryptoSystem interface {
	Crypto(id t.NodeID) (mirCrypto.Crypto, error)
	Module(id t.NodeID) (*modules.SimpleEventApplier, error)
}

type localPseudoCryptoSystem struct {
	nodeIDs       []t.NodeID
	cryptoType    string
	localKeyPairs mirCrypto.KeyPairs
}

// NewLocalCryptoSystem creates an instance of LocalCryptoSystem suitable for tests.
// In the current implementation, cryptoType can only be "pseudo" or "dummy".
// TODO: Think about merging all the code in this file with the pseudo crypto implementation.
func NewLocalCryptoSystem(cryptoType string, nodeIDs []t.NodeID, _ logging.Logger) (LocalCryptoSystem, error) {
	// Generate keys once. All crypto modules derived from this crypto system will use those keys.
	keyPairs, err := mirCrypto.GenerateKeys(len(nodeIDs), mirCrypto.DefaultPseudoSeed)
	if err != nil {
		return nil, err
	}

	return &localPseudoCryptoSystem{nodeIDs, cryptoType, keyPairs}, nil
}

func (cs *localPseudoCryptoSystem) Crypto(id t.NodeID) (mirCrypto.Crypto, error) {
	if cs.cryptoType == "pseudo" || cs.cryptoType == "pseudo-purego" {
		return mirCrypto.InsecureCryptoForTestingOnly(cs.nodeIDs, id, &cs.localKeyPairs)
	} else if cs.cryptoType == "dummy" {
		return &mirCrypto.DummyCrypto{
			DummySig: []byte("a valid signature"),
		}, nil
	} else {
		return nil, es.Errorf("unknown local crypto system type: %v (must be pseudo or dummy)", cs.cryptoType)
	}
}

func (cs *localPseudoCryptoSystem) Module(id t.NodeID) (*modules.SimpleEventApplier, error) {
	c, err := cs.Crypto(id)
	if err != nil {
		return nil, err
	}
	return mirCrypto.New(c), nil
}
