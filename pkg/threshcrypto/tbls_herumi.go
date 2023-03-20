package threshcrypto

import (
	"fmt"
	"io"

	"github.com/herumi/bls-eth-go-binary/bls"
	"golang.org/x/exp/slices"

	t "github.com/filecoin-project/mir/pkg/types"
)

func init() {
	if err := bls.Init(bls.BLS12_381); err != nil {
		panic(fmt.Errorf("failed to init herumi bls lib: %w", err))
	}

	sk := new(bls.SecretKey)
	sk.SetByCSPRNG()

	pk, err := sk.GetSafePublicKey()
	if err != nil {
		panic(err)
	}

	msg := []byte{1, 2, 3, 42}

	sig := sk.SignByte(msg)

	if !sig.VerifyByte(pk, msg) {
		panic("herumi broken")
	}

	ser := sig.Serialize()
	sig2 := new(bls.Sign)
	if err := sig2.Deserialize(ser); err != nil {
		panic(err)
	}

	if !sig2.VerifyByte(pk, msg) {
		panic("herumi broken")
	}
}

// HerumiTBLSInst an instance of a BLS-based (t, len(members))-threshold signature scheme
// It is capable of creating signature shares with its (single) private key share,
// and validating/recovering signatures involving all group members.
type HerumiTBLSInst struct {
	ownIdx  int
	T       int
	members []t.NodeID

	skShare  *bls.SecretKey
	pk       *bls.PublicKey
	pkShares []*bls.PublicKey
}

// HerumiTBLSKeygen constructs a set BLSTInst for a given set of member nodes and threshold T
// using the BLS12-381 pairing, with signatures being points on curve G1, and keys points on curve G2.
func HerumiTBLSKeygen(T int, members []t.NodeID, randSource io.Reader) ([]*HerumiTBLSInst, error) {
	bls.SetRandFunc(randSource)

	instances := make([]*HerumiTBLSInst, len(members))
	pkShares := make([]*bls.PublicKey, len(members))

	masterSk := new(bls.SecretKey)
	masterSk.SetByCSPRNG()

	skComponents := masterSk.GetMasterSecretKey(T)
	pk, err := skComponents[0].GetSafePublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to compute group public key: %w", err)
	}

	for i := range instances {
		instances[i] = &HerumiTBLSInst{
			ownIdx:  i,
			T:       T,
			members: members,
		}

		id, err := herumiID(uint64(i))
		if err != nil {
			return nil, fmt.Errorf("failed to generate herumi ID: %w", err)
		}

		skShare := new(bls.SecretKey)
		if err := skShare.Set(skComponents, &id); err != nil {
			return nil, fmt.Errorf("failed to create secret key share #%d: %w", i, err)
		}

		pkShare, err := skShare.GetSafePublicKey()
		if err != nil {
			return nil, fmt.Errorf("failed to compute group public key: %w", err)
		}

		instances[i].skShare = skShare
		pkShares[i] = pkShare
	}

	// save public key and shares
	for i := range instances {
		instances[i].pk = pk
		instances[i].pkShares = pkShares
	}

	return instances, nil
}

// SignShare constructs a signature share for the message.
func (inst *HerumiTBLSInst) SignShare(msg [][]byte) ([]byte, error) {
	sig := inst.skShare.SignByte(digest(msg))
	return serializeHerumiSigShare(sig, uint64(inst.ownIdx)), nil
}

// VerifyShare verifies that a signature share is for a given message from a given node.
func (inst *HerumiTBLSInst) VerifyShare(msg [][]byte, sigShare []byte, nodeID t.NodeID) error {
	presumedID := slices.Index(inst.members, nodeID)
	if presumedID == -1 {
		return fmt.Errorf("invalid signer: %v", nodeID)
	}

	sig, id, err := deserializeHerumiSigShare(sigShare)
	if err != nil {
		return fmt.Errorf("failed to parse share: %w", err)
	}

	if id.GetDecString() != fmt.Sprintf("%d", presumedID+1) {
		return fmt.Errorf("sig owner mismatch")
	}

	if sig.VerifyByte(inst.pkShares[presumedID], digest(msg)) {
		return nil
	}

	return fmt.Errorf("verification failed")
}

// VerifyFull verifies that a (full) signature is valid for a given message.
func (inst *HerumiTBLSInst) VerifyFull(msg [][]byte, sigFull []byte) error {
	sig := bls.Sign{}
	if err := sig.Deserialize(sigFull); err != nil {
		return fmt.Errorf("error deserializing sig: %w", err)
	}

	if sig.VerifyByte(inst.pk, digest(msg)) {
		return nil
	}

	return fmt.Errorf("verification failed")
}

// Recover recovers a full signature from a set of (previously validated) shares, that are known to be from
// distinct nodes.
func (inst *HerumiTBLSInst) Recover(msg [][]byte, sigShares [][]byte) ([]byte, error) {
	fullSig := bls.Sign{}

	if len(sigShares) < inst.T {
		// herumi BLS recover operation will succeed, but yield a bad signature
		return nil, fmt.Errorf("not enough sig shares")
	}

	parsedShares := make([]bls.Sign, len(sigShares))
	parsedIDs := make([]bls.ID, len(sigShares))
	for i, share := range sigShares {
		parsedShare, parsedID, err := deserializeHerumiSigShare(share)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize share #%d: %w", i, err)
		}

		parsedShares[i] = *parsedShare
		parsedIDs[i] = parsedID
	}

	if err := fullSig.Recover(parsedShares, parsedIDs); err != nil {
		return nil, fmt.Errorf("failed to recover full sig: %w", err)
	}

	return fullSig.Serialize(), nil
}

func serializeHerumiSigShare(sig *bls.Sign, id uint64) []byte {
	sigBytes := sig.Serialize()

	serialized := make([]byte, 0, len(sigBytes)+8)

	serialized = append(serialized, t.Uint64ToBytes(id)...)
	if len(serialized) != 8 {
		panic("uint64 are supposed to be 8 bytes long")
	}

	serialized = append(serialized, sigBytes...)

	return serialized
}

func deserializeHerumiSigShare(serialized []byte) (*bls.Sign, bls.ID, error) {
	idNum := t.Uint64FromBytes(serialized[0:8])
	id, err := herumiID(idNum)
	if err != nil {
		return nil, id, err
	}

	sig := &bls.Sign{}
	if err := sig.Deserialize(serialized[8:]); err != nil {
		return sig, id, err
	}

	return sig, id, nil
}

func herumiID(idNum uint64) (bls.ID, error) {
	id := bls.ID{}

	if err := id.SetDecString(fmt.Sprintf("%d", idNum+1)); err != nil {
		return id, fmt.Errorf("failed to create key ID: %w", err)
	}

	return id, nil
}
