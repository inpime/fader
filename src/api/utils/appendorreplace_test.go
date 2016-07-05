package utils

import (
	"reflect"
	"testing"
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
	}

	AppendOrReplace(&d, s)

	if !reflect.DeepEqual(d, e) {
		t.FailNow()
	}
}
