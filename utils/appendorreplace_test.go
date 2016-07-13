package utils

import (
	// "github.com/inpime/fader/utils/sdata"
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
	// M *sdata.StringMap
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
		// sdata.NewStringMapFrom(map[string]interface{}{
		// 	"a":  "b",
		// 	"f1": sdata.NewStringMapFrom(map[string]interface{}{"a": "b"}),
		// }),
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
		// sdata.NewStringMapFrom(map[string]interface{}{
		// 	"c":  "d",
		// 	"f1": sdata.NewStringMapFrom(map[string]interface{}{"c": "d"}),
		// 	"f2": sdata.NewStringMapFrom(map[string]interface{}{"a": "b"}),
		// }),
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
		// sdata.NewStringMapFrom(map[string]interface{}{
		// 	"a":  "b",
		// 	"c":  "d",
		// 	"f1": sdata.NewStringMapFrom(map[string]interface{}{"a": "b"}),
		// 	"f2": sdata.NewStringMapFrom(map[string]interface{}{"a": "b"}),
		// }),
	}

	AppendOrReplace(&d, s)

	t.Logf("%#v", d)

	if !reflect.DeepEqual(d, e) {
		t.FailNow()
	}
}
