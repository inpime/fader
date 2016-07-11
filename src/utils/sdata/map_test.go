package sdata

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/vmihailenco/msgpack.v2"
	"testing"
)

func encode(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func decodeStr(v interface{}, b []byte) error {
	dec := msgpack.NewDecoder(bytes.NewBuffer(b))
	dec.DecodeMapFunc = func(d *msgpack.Decoder) (interface{}, error) {
		n, err := d.DecodeMapLen()
		if err != nil {
			return nil, err
		}

		m := make(map[string]interface{}, n)
		for i := 0; i < n; i++ {
			mk, err := d.DecodeString()
			if err != nil {
				return nil, err
			}

			mv, err := d.DecodeInterface()
			if err != nil {
				return nil, err
			}

			m[mk] = mv
		}
		return m, nil
	}

	return dec.Decode(v)
}

func decode(v interface{}, b []byte) error {
	return msgpack.Unmarshal(b, v)
}

type testStruct struct {
	A string
	B float64
}

func checkMap(m *Map, t *testing.T) {
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

func TestMap_load(t *testing.T) {
	m := NewMap()
	m.LoadFrom(map[interface{}]interface{}{
		"a":   "b",
		"c":   1.23,
		"d":   2,
		"arr": []interface{}{"a", "b", "c"},
		"mapstr": map[interface{}]interface{}{
			"a":   "b",
			"c":   1.23,
			"d":   2,
			"arr": []interface{}{"a", "b", "c"},
			"mapstr": map[interface{}]interface{}{
				"a": "b",
				"c": 1.23,
				"d": 2,
			},
		},
	})
	checkMap(m, t)
}

func TestMap_simple(t *testing.T) {
	msgpack.RegisterExt(1, testStruct{})

	m := NewMap()
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

	m = NewMap()
	err = decode(m, networkBytes)
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
}
