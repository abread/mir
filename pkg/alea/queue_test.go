package alea

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfVec_Mutation(t *testing.T) {
	v := NewInfVec[string](3)

	el0, ok0 := v.Get(0)
	assert.True(t, ok0)

	*el0 = "asd"

	el0Again, ok0Again := v.Get(0)
	assert.True(t, ok0Again)
	assert.Equal(t, "asd", *el0Again)
}

func TestInfVec_GetNonFirst(t *testing.T) {
	v := NewInfVec[string](3)

	el1, ok1 := v.Get(1)

	assert.True(t, ok1)
	assert.Equal(t, "", *el1)
}

func TestInfVec_Recycling(t *testing.T) {
	v := NewInfVec[string](3)

	el0, ok0 := v.Get(0)
	assert.True(t, ok0)
	*el0 = "asdef"
	v.Get(1)
	v.Get(2)

	_, ok3PreFree := v.Get(3)
	assert.False(t, ok3PreFree)

	v.Free(0)

	_, ok0Freed := v.Get(0)
	assert.False(t, ok0Freed)

	el3, ok3 := v.Get(3)
	assert.True(t, ok3)
	assert.Equal(t, "", *el3)
}
