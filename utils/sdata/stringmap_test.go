package sdata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func checkStringMap(m *StringMap, t *testing.T) {
	assert.Equal(t, m.Size(), 5)
	assert.Equal(t, m.String("a"), "b")
	assert.Equal(t, m.Float("c"), 1.23)
	assert.Equal(t, m.Int("d"), 2)
	assert.Equal(t, m.A("arr").Size(), 3)
	assert.Equal(t, m.A("arr").Includes("a", "b", "c"), true)
	assert.Equal(t, m.M("mapstr").String("a"), "b")
	assert.Equal(t, m.M("mapstr").Float("c"), 1.23)
	assert.Equal(t, m.M("mapstr").Int("d"), 2)
	assert.Equal(t, m.M("mapstr").A("arr").Includes("a", "b", "c"), true)
	assert.Equal(t, m.M("mapstr").A("arr").Size(), 3)
	assert.Equal(t, m.M("mapstr").M("mapstr").String("a"), "b")
	assert.Equal(t, m.M("mapstr").M("mapstr").Float("c"), 1.23)
	assert.Equal(t, m.M("mapstr").M("mapstr").Int("d"), 2)
}

func TestStringMap_load(t *testing.T) {
	m := NewStringMap()
	m.LoadFrom(map[string]interface{}{
		"a":   "b",
		"c":   1.23,
		"d":   2,
		"arr": []interface{}{"a", "b", "c"},
		"mapstr": map[string]interface{}{
			"a":   "b",
			"c":   1.23,
			"d":   2,
			"arr": []interface{}{"a", "b", "c"},
			"mapstr": map[string]interface{}{
				"a": "b",
				"c": 1.23,
				"d": 2,
			},
		},
	})
	checkStringMap(m, t)

	m = NewStringMap()
	m.LoadFrom(`{
        "a": "b",
        "c": 1.23,
        "d": 2,
        "arr": ["a", "b", "c"],
        "mapstr": {
            "a": "b",
            "c": 1.23,
            "d": 2,
            "arr": ["a", "b", "c"],
            "mapstr": {
                "a": "b",
                "c": 1.23,
                "d": 2
            }
        }
    }`)
	checkStringMap(m, t)
}

func TestStringMap_simple(t *testing.T) {

	m := NewStringMap()
	m.Set("s1", "v1")
	m.Set("s2", "v2")
	m.Set("i1", 2)
	m.Set("f1", 1.2)
	m.M("map1").Set("s1", "v1")
	m.M("map1").Set("b1", true)
	m.M("map1").Set("b2", false)
	m.A("arr1").Add("a").Add("b", "c").Add(1.2)
	m.Set("struct", testStruct{A: "a", B: 1.23})

	assert.Equal(t, m.Size(), 7)
	assert.Equal(t, m.Int("i1"), 2)
	assert.Equal(t, m.Int64("i1"), int64(2))
	assert.Equal(t, m.Float("i1"), float64(2))
	assert.Equal(t, m.Float("f1"), 1.2)
	assert.Equal(t, m.M("map1").String("s1"), "v1")
	assert.Equal(t, m.M("map1").Bool("b1"), true)
	assert.Equal(t, m.M("map1").Bool("b2"), false)
	assert.Equal(t, m.A("arr1").Size(), 4)
	assert.Equal(t, m.A("arr1").Includes(1.2, "a", "b", "c"), true)
	assert.Equal(t, m.GetOrNil("struct").(testStruct).A, "a")
	assert.Equal(t, m.GetOrNil("struct").(testStruct).B, 1.23)

	networkBytes, err := encode(m)
	assert.NoError(t, err)

	// bytes

	m = NewStringMap()
	err = decodeStr(m, networkBytes)
	assert.NoError(t, err)

	assert.Equal(t, m.Size(), 7)
	assert.Equal(t, m.Int("i1"), 2)
	assert.Equal(t, m.Int64("i1"), int64(2))
	assert.Equal(t, m.Float("i1"), float64(2))
	assert.Equal(t, m.Float("f1"), 1.2)
	assert.Equal(t, m.M("map1").String("s1"), "v1")
	assert.Equal(t, m.M("map1").Bool("b1"), true)
	assert.Equal(t, m.M("map1").Bool("b2"), false)
	assert.Equal(t, m.A("arr1").Size(), 4)
	assert.Equal(t, m.A("arr1").Includes(1.2, "a", "b", "c"), true)
	assert.Equal(t, m.GetOrNil("struct").(testStruct).A, "a")
	assert.Equal(t, m.GetOrNil("struct").(testStruct).B, 1.23)

	// changes from pointer

	_m := map[string]interface{}{
		"a":   "b",
		"arr": []interface{}{"a", "b", "c"},
	}
	m = NewStringMapFrom(_m)
	m.Set("a", "d")
	m.A("arr").Remove("a")

	assert.Equal(t, m.String("a"), "d")
	assert.Equal(t, m.A("arr").Size(), 2)
	assert.Equal(t, m.A("arr").Includes("b", "c"), true)
	assert.Equal(t, _m["a"], "d")
	assert.Equal(t, NewArrayFrom(_m["arr"]).Size(), 2)
	assert.Equal(t, NewArrayFrom(_m["arr"]).Includes("b", "c"), true)
}
