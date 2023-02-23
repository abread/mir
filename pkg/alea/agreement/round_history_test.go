package agreement

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertAbsent(t assert.TestingT, rh *AgRoundHistory, idx uint64) {
	_, ok := rh.Get(idx)
	assert.False(t, ok)
}
func assertPresent(t assert.TestingT, rh *AgRoundHistory, idx uint64, expected bool) {
	val, ok := rh.Get(idx)
	assert.True(t, ok)
	assert.Equal(t, expected, val)
}

func TestRoundHistory(t *testing.T) {
	var rh AgRoundHistory

	assert.Equal(t, uint64(0), rh.Len(), "empty at first")
	assertAbsent(t, &rh, 0)
	assertAbsent(t, &rh, subSliceSize-1)
	assertAbsent(t, &rh, subSliceSize)
	assertAbsent(t, &rh, subSliceSize+1)

	rh.Push(false)
	assert.Equal(t, uint64(1), rh.Len(), "pushed one")
	assertPresent(t, &rh, 0, false)
	assertAbsent(t, &rh, subSliceSize-1)
	assertAbsent(t, &rh, subSliceSize)
	assertAbsent(t, &rh, subSliceSize+1)

	rh.Push(true)
	assertPresent(t, &rh, 0, false)
	assertPresent(t, &rh, 1, true)

	assert.Equal(t, 1, len(rh))

	for i := 2; i < subSliceSize; i++ {
		rh.Push(false)
	}
	for i := uint64(2); i < uint64(subSliceSize); i++ {
		assertPresent(t, &rh, i, false)
	}
	assert.Equal(t, 1, len(rh))
	assert.Equal(t, false, rh.hasFreeSpace())

	for i := 0; i < subSliceSize; i++ {
		rh.Push(true)
	}
	for i := uint64(2); i < uint64(subSliceSize); i++ {
		assertPresent(t, &rh, i, false)
		assertPresent(t, &rh, i+uint64(subSliceSize), true)
	}
	assert.Equal(t, 2, len(rh))
	assert.Equal(t, false, rh.hasFreeSpace())
}
