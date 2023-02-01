package modring

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWSC_PutMutateFree(t *testing.T) {
	windowController := NewRingController(8)

	// window should be bounded
	assert.False(t, windowController.TryAcquire(8))

	for _, i := range []uint64{0, 7, 1, 2, 3, 4, 5, 6} { // non-monotonic orders should work
		assert.False(t, windowController.IsSlotInUse(i))

		assert.True(t, windowController.TryAcquire(i))
	}

	for i := uint64(0); i < 8; i++ {
		assert.True(t, windowController.IsSlotInUse(i))
	}

	// queue should be bounded
	assert.False(t, windowController.TryAcquire(8))

	// after freeing the first element, the window should advance
	assert.True(t, windowController.TryFree(0))

	assert.True(t, windowController.TryAcquire(8))
}

func TestWSC_OutOfOrderViewChange(t *testing.T) {
	windowController := NewRingController(3)

	for i := uint64(0); i < 3; i++ {
		require.True(t, windowController.TryAcquire(i))
	}

	// slot 3-5 currently out of view
	for i := uint64(3); i < 6; i++ {
		assert.False(t, windowController.TryAcquire(i))
	}

	// even when freeing slots 1 and 2, 3-5 are out of view
	require.True(t, windowController.TryFree(1))
	require.True(t, windowController.TryFree(2))
	for i := uint64(3); i < 6; i++ {
		assert.False(t, windowController.TryAcquire(i))
	}

	// after freeing slot 0, slots 3-5 become available
	require.True(t, windowController.TryFree(0))
	for i := uint64(3); i < 6; i++ {
		assert.True(t, windowController.TryAcquire(i))
	}
}
