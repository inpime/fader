package sdata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArray_simple(t *testing.T) {
	a := NewArray()
	a.Add("a")
	a.Add("b", "c").Add("d")
	a.Remove("a")
	a.Add("e")

	assert.Equal(t, a.Size(), 4)
	assert.Equal(t, a.Includes("a"), false)
	assert.Equal(t, a.Includes("e", "b", "c", "d"), true)

	network, err := encode(a)
	assert.NoError(t, err)

	//

	a = NewArray()
	err = decode(a, network)
	assert.NoError(t, err)

	assert.Equal(t, a.Size(), 4)
	assert.Equal(t, a.Includes("a"), false)
	assert.Equal(t, a.Includes("e", "b", "c", "d"), true)
}
