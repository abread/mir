package threshcrypto

import (
	"bytes"

	es "github.com/go-errors/errors"

	t "github.com/filecoin-project/mir/pkg/types"
)

// DummyCrypto represents a dummy MirModule module that
// always produces a SHA256 hash of the data as the full signature.
// Signature shares always consist of the nodeID followed by a preset suffix (DummySigShareSuffix)
// Verification of these dummy signatures always succeeds.
// This is intended as a stub for testing purposes.
type DummyCrypto struct {
	// The only accepted signature share suffix
	DummySigShareSuffix []byte

	// Current node ID
	NodeID t.NodeID
}

// SignShare always returns the dummy signature DummySig, regardless of the data.
func (dc *DummyCrypto) SignShare(_ [][]byte) ([]byte, error) {
	return dc.buildSigShare(dc.NodeID), nil
}

// VerifyShare returns nil (i.e. success) only if signature share equals nodeID||DummySigShareSuffix.
// data is ignored.
func (dc *DummyCrypto) VerifyShare(_ [][]byte, sigShare []byte, nodeID t.NodeID) error {
	if !bytes.Equal(sigShare, dc.buildSigShare(nodeID)) {
		return es.Errorf("dummy signature mismatch")
	}

	return nil
}

// VerifyFull returns nil (i.e. success) only if signature equals DummySig.
// data is ignored.
func (dc *DummyCrypto) VerifyFull(data [][]byte, signature []byte) error {
	expectedSig := digest(data)

	if !bytes.Equal(signature, expectedSig) {
		return es.Errorf("dummy signature mismatch")
	}

	return nil
}

// Recovers full signature from signature shares if they are valid, otherwise an error is returned.
// data is ignored.
func (dc *DummyCrypto) Recover(data [][]byte, sigShares [][]byte) ([]byte, error) {
	for _, share := range sigShares {
		nodeID := share[:(len(share) - len(dc.DummySigShareSuffix))]

		if err := dc.VerifyShare(data, share, t.NodeID(nodeID)); err != nil {
			return nil, err
		}
	}

	return digest(data), nil
}

// construct the dummy signature share for a nodeID
func (dc *DummyCrypto) buildSigShare(nodeID t.NodeID) []byte {
	share := []byte(nodeID)
	share = append(share, dc.DummySigShareSuffix...)
	return share
}
