package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func mBasicStrategy(t *testing.T, v map[string]interface{}) {

	m := Map(v)

	m.Set("s1", "v1")
	m.Set("s2", "v2")
	m.Set("i1", 1)
	m.Set("f1", 1.2)
	m.M("map1").Set("s1", "v1")
	m.M("map1").Set("b1", true)
	m.M("map1").Set("b2", false)

	if t != nil {
		assert.Equal(t, m.String("s1"), "v1")
		assert.Equal(t, m.String("s2"), "v2")
		assert.Equal(t, m.Int("i1"), 1)
		assert.Equal(t, m.Int64("i1"), int64(1))
		assert.Equal(t, m.Float64("i1"), float64(1))
		assert.Equal(t, m.Float64("f1"), 1.2)
		assert.Equal(t, m.M("map1").String("s1"), "v1")
		assert.Equal(t, m.M("map1").Bool("b1"), true)
		assert.Equal(t, m.M("map1").Bool("b2"), false)
	}
}

func TestM_basic(t *testing.T) {
	v := map[string]interface{}{}
	mBasicStrategy(t, v)
}

func BenchmarkM_basic(b *testing.B) {
	v := map[string]interface{}{}
	for i := 0; i < b.N; i++ {
		mBasicStrategy(nil, v)
	}
}
