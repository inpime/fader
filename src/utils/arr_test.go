package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestA_basic(t *testing.T) {
	a := NewA([]string{"v"})
	a.Add("a").Add("b").Delete("a")
	a.Delete("c").Add("c")

	assert.Equal(t, a.Include("a"), false)
	assert.Equal(t, a.Include("b"), true)
	assert.Equal(t, a.Include("c"), true)
	assert.Equal(t, a.Include("v"), true)
	assert.Equal(t, a.Len(), 3)
}
