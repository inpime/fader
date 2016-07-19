package utils

import (
	"reflect"
	"testing"
)

type simplestruct struct {
	A string
	B string
}

func TestMerge_simplestruct(t *testing.T) {
	src := simplestruct{
		A: "a",
		B: "b",
	}

	dst := simplestruct{
		A: "a",
		B: "c",
	}

	exp := simplestruct{
		A: "a",
		B: "b",
	}

	Merge(&dst, src)

	if !reflect.DeepEqual(dst, exp) {
		t.FailNow()
	}
}

func TestMerge_simplemap(t *testing.T) {
	src := map[string]interface{}{
		"a": "a",
		"b": "b",
	}

	dst := map[string]interface{}{
		"a": "d",
		"b": "e",
		"c": "c",
	}

	exp := map[string]interface{}{
		"a": "a",
		"b": "b",
		"c": "c",
	}

	Merge(&dst, src)

	if !reflect.DeepEqual(dst, exp) {
		t.FailNow()
	}
}

type structmap struct {
	A map[string]interface{}
	B string
}

func TestMerge_structmap(t *testing.T) {
	src := structmap{
		A: map[string]interface{}{
			"a": "a",
			"b": "b",
			"d": "d",
		},
		B: "b",
	}

	dst := structmap{
		A: map[string]interface{}{
			"a": "d",
			"b": "e",
			"c": "c",
		},
		B: "c",
	}

	exp := structmap{
		A: map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
			"d": "d",
		},
		B: "b",
	}

	Merge(&dst, src)

	if !reflect.DeepEqual(dst, exp) {
		t.FailNow()
	}
}

type structmapstruct struct {
	A map[string]simplestruct
	B string
}

func TestMerge_structmapstruct(t *testing.T) {
	src := structmapstruct{
		A: map[string]simplestruct{
			"a": simplestruct{
				A: "aa",
				B: "ab",
			},
			"b": simplestruct{
				A: "ba",
				B: "bb",
			},
			"d": simplestruct{
				A: "da",
				B: "db",
			},
		},
		B: "b",
	}

	dst := structmapstruct{
		A: map[string]simplestruct{
			"a": simplestruct{
				A: "not valid",
				B: "not valid",
			},
			"b": simplestruct{
				A: "not valid",
				B: "not valid",
			},
			"c": simplestruct{
				A: "ca",
				B: "cb",
			},
		},
		B: "c",
	}

	exp := structmapstruct{
		A: map[string]simplestruct{
			"a": simplestruct{
				A: "aa",
				B: "ab",
			},
			"b": simplestruct{
				A: "ba",
				B: "bb",
			},
			"c": simplestruct{
				A: "ca",
				B: "cb",
			},
			"d": simplestruct{
				A: "da",
				B: "db",
			},
		},
		B: "b",
	}

	Merge(&dst, src)

	if !reflect.DeepEqual(dst, exp) {
		t.FailNow()
	}
}

func TestMerge_mapinterface(t *testing.T) {
	src := map[string]interface{}{
		"a": map[string]interface{}{
			"b": "b",
		},
		"b": map[string]interface{}{
			"a": "a",
			"b": "b",
		},
		"arr":  []string{"a", "b"},
		"arr1": []string{"a", "b"},
		"arr3": []simplestruct{
			simplestruct{"1", "2"},
			simplestruct{"3", "4"},
		},
		"arr4": []map[string]interface{}{
			map[string]interface{}{
				"1": "2",
			},
			map[string]interface{}{
				"3": "4",
			},
		},
	}

	dst := map[string]interface{}{
		"a": map[string]interface{}{
			"a": "a",
			"b": "not valid",
		},
		"b": map[string]interface{}{
			"a": "not valid",
			"b": "not valid",
		},
		"c": map[string]interface{}{
			"a": "a",
			"b": "b",
		},
		"arr":  []string{"c", "d"},
		"arr2": []string{"a", "b"},
		"arr3": []simplestruct{
			simplestruct{"5", "6"},
			simplestruct{"7", "8"},
		},
		"arr4": []map[string]interface{}{
			map[string]interface{}{
				"5": "6",
			},
			map[string]interface{}{
				"7": "8",
			},
		},
	}

	exp := map[string]interface{}{
		"a": map[string]interface{}{
			"a": "a",
			"b": "b",
		},
		"b": map[string]interface{}{
			"a": "a",
			"b": "b",
		},
		"c": map[string]interface{}{
			"a": "a",
			"b": "b",
		},
		"arr":  []string{"c", "d", "a", "b"},
		"arr1": []string{"a", "b"},
		"arr2": []string{"a", "b"},
		"arr3": []simplestruct{
			// important the order of
			simplestruct{"5", "6"},
			simplestruct{"7", "8"},
			simplestruct{"1", "2"},
			simplestruct{"3", "4"},
		},
		"arr4": []map[string]interface{}{
			// important the order of
			map[string]interface{}{
				"5": "6",
			},
			map[string]interface{}{
				"7": "8",
			},
			map[string]interface{}{
				"1": "2",
			},
			map[string]interface{}{
				"3": "4",
			},
		},
	}

	Merge(&dst, src)

	if !reflect.DeepEqual(dst, exp) {
		t.FailNow()
	}
}
