package utils

import (
	"reflect"
	"testing"
	"utils"
)

type simpleSubType struct {
	A string
}

type simpleType struct {
	A []string
	B string
	C []simpleSubType
	D map[string]string
	E map[string]simpleSubType
	G *simpleSubType
	F map[string]*simpleSubType
	M utils.M
}

func TestAppendOrReplace(t *testing.T) {
	s := simpleType{[]string{"a", "b"},
		"a",
		[]simpleSubType{simpleSubType{"a"}, simpleSubType{"b"}},
		map[string]string{"a": "b", "c": "d"},
		map[string]simpleSubType{
			"a": simpleSubType{"b"},
			"c": simpleSubType{"d"},
		},
		&simpleSubType{"a"},
		map[string]*simpleSubType{
			"a": &simpleSubType{"b"},
			"c": &simpleSubType{"d"},
		},
		utils.M(map[string]interface{}{
			"a":  "b",
			"f1": utils.M(map[string]interface{}{"a": "b"}),
		}),
	}
	d := simpleType{[]string{"b", "c"},
		"b",
		[]simpleSubType{simpleSubType{"c"}, simpleSubType{"d"}},
		map[string]string{"a": "z", "c": "y", "e": "f"},
		map[string]simpleSubType{
			"a": simpleSubType{"z"},
			"c": simpleSubType{"y"},
			"e": simpleSubType{"f"},
		},
		&simpleSubType{"b"},
		map[string]*simpleSubType{
			"a": &simpleSubType{"z"},
			"c": &simpleSubType{"y"},
			"e": &simpleSubType{"f"},
		},
		utils.M(map[string]interface{}{
			"c":  "d",
			"f1": utils.M(map[string]interface{}{"c": "d"}),
			"f2": utils.M(map[string]interface{}{"a": "b"}),
		}),
	}

	// expected
	e := simpleType{
		[]string{"b", "c", "a", "b"},
		"a",
		[]simpleSubType{simpleSubType{"c"}, simpleSubType{"d"}, simpleSubType{"a"}, simpleSubType{"b"}},
		map[string]string{"a": "b", "c": "d", "e": "f"},
		map[string]simpleSubType{
			"a": simpleSubType{"b"},
			"c": simpleSubType{"d"},
			"e": simpleSubType{"f"},
		},
		&simpleSubType{"a"},
		map[string]*simpleSubType{
			"a": &simpleSubType{"b"},
			"c": &simpleSubType{"d"},
			"e": &simpleSubType{"f"},
		},
		utils.M(map[string]interface{}{
			"a":  "b",
			"c":  "d",
			"f1": utils.M(map[string]interface{}{"a": "b"}),
			"f2": utils.M(map[string]interface{}{"a": "b"}),
		}),
	}

	AppendOrReplace(&d, s)

	if !reflect.DeepEqual(d, e) {
		t.FailNow()
	}
}
