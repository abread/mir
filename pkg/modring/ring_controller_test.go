package modring

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWSC_PutMutateFree(t *testing.T) {
	windowController := NewRingController(8)

	// window should be bounded
	assert.False(t, windowController.MarkCurrentIfFuture(8))

	for _, i := range []uint64{0, 7, 1, 2, 3, 4, 5, 6} { // non-monotonic orders should work
		assert.False(t, windowController.IsCurrentSlot(i))

		assert.True(t, windowController.MarkCurrentIfFuture(i))
	}

	for i := uint64(0); i < 8; i++ {
		assert.True(t, windowController.IsCurrentSlot(i))
	}

	// queue should be bounded
	assert.False(t, windowController.MarkCurrentIfFuture(8))

	// after freeing the first element, the window should advance
	assert.NoError(t, windowController.MarkPast(0))

	assert.True(t, windowController.MarkCurrentIfFuture(8))
}

func TestWSC_OutOfOrderViewChange(t *testing.T) {
	windowController := NewRingController(3)

	for i := uint64(0); i < 3; i++ {
		require.True(t, windowController.MarkCurrentIfFuture(i))
	}

	// slot 3-5 currently out of view
	for i := uint64(3); i < 6; i++ {
		assert.False(t, windowController.MarkCurrentIfFuture(i))
	}

	// even when freeing slots 1 and 2, 3-5 are out of view
	require.NoError(t, windowController.MarkPast(1))
	require.NoError(t, windowController.MarkPast(2))
	for i := uint64(3); i < 6; i++ {
		assert.False(t, windowController.MarkCurrentIfFuture(i))
	}

	// after freeing slot 0, slots 3-5 become available
	require.NoError(t, windowController.MarkPast(0))
	for i := uint64(3); i < 6; i++ {
		assert.True(t, windowController.MarkCurrentIfFuture(i))
	}
}

func TestWSC_Load(t *testing.T) {
	windowController := NewRingController(2)

	for i := uint64(0); i < 2; i++ {
		require.True(t, windowController.MarkCurrentIfFuture(i))
	}

	for i := uint64(2); i < 64; i++ {
		assert.False(t, windowController.MarkCurrentIfFuture(i))
		assert.NoError(t, windowController.MarkPast(i-2))
		assert.True(t, windowController.MarkCurrentIfFuture(i))
	}

	assert.False(t, windowController.MarkCurrentIfFuture(64))
	assert.NoError(t, windowController.MarkPast(64-1))
	assert.False(t, windowController.MarkCurrentIfFuture(64))
	assert.NoError(t, windowController.MarkPast(64-2))
	assert.True(t, windowController.MarkCurrentIfFuture(64))
}
