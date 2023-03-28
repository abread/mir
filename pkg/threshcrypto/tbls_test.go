package threshcrypto

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/types"
)

type TestKeygen func(T, N int) []ThreshCrypto

func BenchmarkTBLSSignShare(b *testing.B) {
	benchmarkSignShare(b, keygenTBLS)
}

func BenchmarkHerumiTBLSSignShare(b *testing.B) {
	benchmarkSignShare(b, keygenHerumiTBLS)
}

func benchmarkSignShare(b *testing.B, keygen TestKeygen) {
	F := 33
	N := 3*F + 1
	keys := keygen(2*F+1, N)
	data := [][]byte{{1, 2, 3, 4, 5}, {4, 2}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := keys[i%N].SignShare(data)
		assert.Nil(b, err)
	}
}

func BenchmarkTBLSVerifyShare(b *testing.B) {
	benchmarkVerifyShare(b, keygenTBLS)
}

func BenchmarkHerumiTBLSVerifyShare(b *testing.B) {
	benchmarkVerifyShare(b, keygenHerumiTBLS)
}

func benchmarkVerifyShare(b *testing.B, keygen TestKeygen) {
	F := 33
	N := 3*F + 1
	keys := keygen(2*F+1, N)
	data := [][]byte{{1, 2, 3, 4, 5}, {4, 2}}

	shares := make([][]byte, N)
	for i := range shares {
		share, err := keys[i].SignShare(data)
		require.NoError(b, err)
		shares[i] = share
	}

	badShares := make([][]byte, N)
	for i := range badShares {
		badShare := slices.Clone(shares[i])
		badShare[len(badShare)-1] ^= 1

		badShares[i] = badShare
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			err := keys[N-(i%N)-1].VerifyShare(data, shares[i%N], types.NewNodeIDFromInt(i%N))
			assert.Nil(b, err)
		} else {
			err := keys[N-(i%N)-1].VerifyShare(data, badShares[i%N], types.NewNodeIDFromInt(i%N))
			assert.NotNil(b, err)
		}
	}
}

func BenchmarkTBLSRecover(b *testing.B) {
	benchmarkRecover(b, keygenTBLS)
}

func BenchmarkHerumiTBLSRecover(b *testing.B) {
	benchmarkRecover(b, keygenHerumiTBLS)
}

func benchmarkRecover(b *testing.B, keygen TestKeygen) {
	F := 33
	N := 3*F + 1
	keys := keygen(2*F+1, N)
	data := [][]byte{{1, 2, 3, 4, 5}, {4, 2}}

	shares := make([][]byte, N)
	for i := range shares {
		var err error
		shares[i], err = keys[i].SignShare(data)
		require.Nil(b, err)
	}
	shares = append(shares, shares...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sig, err := keys[i%N].Recover(data, shares[i%N:(2*F+1)+i%N])
		assert.Nil(b, err)
		assert.NotNil(b, sig)
	}
}

func BenchmarkTBLSVerifyFull(b *testing.B) {
	benchmarkVerifyFull(b, keygenTBLS)
}

func BenchmarkHerumiTBLSVerifyFull(b *testing.B) {
	benchmarkVerifyFull(b, keygenHerumiTBLS)
}

func benchmarkVerifyFull(b *testing.B, keygen TestKeygen) {
	F := 33
	N := 3*F + 1
	keys := keygen(2*F+1, N)
	data := [][]byte{{1, 2, 3, 4, 5}, {4, 2}}

	shares := make([][]byte, 2*F+1)
	for i := range shares {
		var err error
		shares[i], err = keys[i].SignShare(data)
		require.Nil(b, err)
	}

	sig, err := keys[0].Recover(data, shares)
	require.Nil(b, err)
	require.NotNil(b, sig)

	badSig := slices.Clone(sig)
	badSig[len(badSig)-1] ^= 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			err := keys[i%N].VerifyFull(data, sig)
			assert.Nil(b, err)
		} else {
			err := keys[i%N].VerifyFull(data, badSig)
			assert.NotNil(b, err)
		}
	}
}

func TestTBLSHappySmoke(t *testing.T) {
	testTCHappySmoke(t, keygenTBLS)
}

func TestHerumiTBLSHappySmoke(t *testing.T) {
	testTCHappySmoke(t, keygenHerumiTBLS)
}

func testTCHappySmoke(t *testing.T, keygen TestKeygen) {
	// confirm basic functionality is there (happy path)
	N := 5
	T := 3

	keys := keygen(T, N)

	data := [][]byte{{1, 2, 3, 4, 5}, {4, 2}}

	shares := make([][]byte, 3)
	for i := range shares {
		sh, err := keys[i].SignShare(data)
		assert.NoError(t, err)

		shares[i] = sh
	}

	// everyone can verify everyone's share
	for _, k := range keys {
		for i, sh := range shares {
			require.NoError(t, k.VerifyShare(data, sh, types.NewNodeIDFromInt(i)))
		}
	}

	for _, kR := range keys {
		// everyone can recover the full signature
		full, err := kR.Recover(data, shares)
		assert.NoError(t, err)

		// everyone can verify the recovered signature
		for _, kV := range keys {
			assert.NoError(t, kV.VerifyFull(data, full))
		}
	}
}

func TestTBLSSadSmoke(t *testing.T) {
	testTCSadSmoke(t, keygenTBLS)
}

func TestHerumiTBLSSadSmoke(t *testing.T) {
	testTCSadSmoke(t, keygenHerumiTBLS)
}

func testTCSadSmoke(t *testing.T, keygen TestKeygen) {
	// confirm some basic problems are detected correctly
	N := 5
	T := 3

	keys := keygen(T, N)

	data := [][]byte{{1, 2, 3, 4, 5}, {4, 2}}

	shares := make([][]byte, N)
	for i := range shares {
		sh, err := keys[i].SignShare(data)
		require.NoError(t, err)

		shares[i] = sh
	}

	// all of the same share is no good
	_, err := keys[0].Recover(data, [][]byte{shares[0], shares[0], shares[0]})
	assert.Error(t, err)

	// too little shares is no good
	_, err = keys[0].Recover(data, shares[:T-1])
	assert.Error(t, err)

	// mangle one of the shares
	shares[1][0] ^= 1

	// mangled share fails verification
	for _, k := range keys {
		assert.Error(t, k.VerifyShare(data, shares[1], "1"))
	}

	// mangled share recovers bad sig
	fullSig, err := keys[0].Recover(data, shares[:3])
	if err == nil {
		err = keys[0].VerifyFull(data, fullSig)
		assert.Error(t, err)
	}
}

func TestTBLSMarshalling(t *testing.T) {
	N := 3
	T := N

	keys := keygenTBLSRaw(T, N)

	keys2 := marshalUnmarshalKeys(t, keys)

	data := [][]byte{{1, 2, 3, 4, 5}, {4, 2}}

	// produce all required shares using both sets of keys
	sigShares := make([][]byte, N)
	sigShares2 := make([][]byte, N)
	for i := range sigShares {
		var err error
		sigShares[i], err = keys[i].SignShare(data)
		assert.NoError(t, err)

		sigShares2[i], err = keys2[i].SignShare(data)
		assert.NoError(t, err)
	}

	// check that the both sets of keys can recover the other's signature
	for i := range sigShares {
		var err error

		_, err = keys[i].Recover(data, sigShares2)
		assert.NoError(t, err)

		_, err = keys2[i].Recover(data, sigShares)
		assert.NoError(t, err)
	}
}

func keygenTBLSRaw(T, N int) []*TBLSInst {
	members := make([]types.NodeID, N)
	for i := range members {
		members[i] = types.NewNodeIDFromInt(i)
	}

	rand := pseudorandomStream(DefaultPseudoSeed)
	return TBLS12381Keygen(T, members, rand)
}

func keygenTBLS(T, N int) []ThreshCrypto {
	return castInsts(keygenTBLSRaw(T, N))
}

func keygenHerumiTBLS(T, N int) []ThreshCrypto {
	members := make([]types.NodeID, N)
	for i := range members {
		members[i] = types.NewNodeIDFromInt(i)
	}

	randSource := rand.New(rand.NewSource(uint64(DefaultPseudoSeed)))

	initialRead := make([]byte, 32)
	n, err := randSource.Read(initialRead)
	if n != 32 && err != nil {
		panic("bad random source")
	}

	insts, err := HerumiTBLSKeygen(T, members, randSource)
	if err != nil {
		panic(fmt.Errorf("error generating keys: %w", err))
	}

	return castInsts(insts)
}

func castInsts[T ThreshCrypto](insts []T) []ThreshCrypto {
	castInsts := make([]ThreshCrypto, len(insts))
	for i, inst := range insts {
		castInsts[i] = inst
	}

	return castInsts
}

func marshalUnmarshalKeys(t *testing.T, src []*TBLSInst) []*TBLSInst {
	res := make([]*TBLSInst, len(src))

	pipeR, pipeW := io.Pipe()

	go func() {
		for i := range src {
			_, err := src[i].MarshalTo(pipeW)
			assert.NoError(t, err)
		}

		pipeW.Close()
	}()

	for i := range res {
		res[i] = &TBLSInst{}
		_, err := res[i].UnmarshalFrom(pipeR)
		assert.NoError(t, err)
	}

	data := make([]byte, 1)
	_, err := pipeR.Read(data)
	assert.ErrorIs(t, err, io.EOF)
	pipeR.Close()
	return res
}
