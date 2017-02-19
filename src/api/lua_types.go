package api

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

func ToValueFromLValue(v lua.LValue) interface{} {
	switch v.Type() {
	case lua.LTNumber:
		return float64(v.(lua.LNumber))
	case lua.LTBool:
		return bool(v.(lua.LBool))
	case lua.LTString:
		return string(v.(lua.LString))
	case lua.LTUserData:
		return v.(*lua.LUserData).Value
	case lua.LTTable:
		tbl := v.(*lua.LTable)
		var keys, vals []interface{}

		isArray := true
		counter := 0
		tbl.ForEach(func(k, v lua.LValue) {
			if k.Type() == lua.LTNumber && int(k.(lua.LNumber)) != counter+1 {
				isArray = false
			}
			if k.Type() != lua.LTNumber {
				isArray = false
			}

			keys = append(keys, k.String())
			vals = append(vals, ToValueFromLValue(v))

			counter++
		})

		if isArray {
			return vals
		}

		_vals := make(map[interface{}]interface{}, counter)
		for i := 0; i < counter; i++ {
			_vals[keys[i]] = vals[i]
		}

		return _vals
	}

	return nil
}

func ToLValueOrNil(v interface{}, L *lua.LState) lua.LValue {
	switch v := v.(type) {
	case int, int8, int32, int64, float32, float64,
		uint, uint8, uint16, uint32, uint64:

		return lua.LNumber(ToFloat64(v))
	case bool:
		return lua.LBool(v)
	case string:
		return lua.LString(v)
	case nil:
		return lua.LNil
	case []int, []int64, []float64, []string, []bool, []interface{}:
		tb := L.NewTable()

		var litems []lua.LValue

		// types
		switch v := v.(type) {
		case []int:
			litems = make([]lua.LValue, len(v))
			for index, item := range v {
				if _v := ToLValueOrNil(item, L); _v != nil {
					litems[index] = _v
				}
			}
		case []int64:
			litems = make([]lua.LValue, len(v))
			for index, item := range v {
				if _v := ToLValueOrNil(item, L); _v != nil {
					litems[index] = _v
				}
			}
		case []float64:
			litems = make([]lua.LValue, len(v))
			for index, item := range v {
				if _v := ToLValueOrNil(item, L); _v != nil {
					litems[index] = _v
				}
			}
		case []bool:
			litems = make([]lua.LValue, len(v))
			for index, item := range v {
				if _v := ToLValueOrNil(item, L); _v != nil {
					litems[index] = _v
				}
			}
		case []string:
			litems = make([]lua.LValue, len(v))
			for index, item := range v {
				if _v := ToLValueOrNil(item, L); _v != nil {
					litems[index] = _v
				}
			}
		case []interface{}:
			litems = make([]lua.LValue, len(v))
			for index, item := range v {
				if _v := ToLValueOrNil(item, L); _v != nil {
					litems[index] = _v
				}
			}
		default:
			log.Printf(
				"[ERR] ToLValueOrNil: slice, not expected type value, got %T",
				v,
			)
		}

		// main

		if len(litems) != 0 {
			for _, item := range litems {
				tb.Append(item)
			}
		}

		return tb

	case map[string]string, map[string]interface{}, map[interface{}]interface{}:
		tb := L.NewTable()

		var keys, values []lua.LValue

		// types
		switch v := v.(type) {
		case map[string]string:
		case map[string]interface{}:
			keys = make([]lua.LValue, len(v))
			values = make([]lua.LValue, len(v))
			var seq = 0
			for key, value := range v {
				keys[seq] = ToLValueOrNil(key, L)
				values[seq] = ToLValueOrNil(value, L)
				seq++
			}
		case map[interface{}]interface{}:
		default:
			log.Printf(
				"[ERR] ToLValueOrNil: not expected type value, got %T",
				v,
			)
		}

		for i := 0; i < len(keys); i++ {
			tb.RawSet(keys[i], values[i])
		}

		return tb
	default:
		log.Printf(
			"[ERR] ToLValueOrNil: not expected type value, got %T",
			v,
		)
	}

	return lua.LNil
}
