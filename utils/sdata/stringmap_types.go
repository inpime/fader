package sdata

import (
	"encoding/json"
)

func (m StringMap) String(key string) string {
	return toString(m.GetOrNil(key))
}

func (m StringMap) Float(key string) float64 {
	return toFloat64(m.GetOrNil(key))
}

func (m StringMap) Int(key string) int {
	return toInt(m.GetOrNil(key))
}

func (m StringMap) Int64(key string) int64 {
	return toInt64(m.GetOrNil(key))
}

func (m StringMap) Bool(key string) bool {
	return toBool(m.GetOrNil(key))
}

func (m *StringMap) Map(key string) map[string]interface{} {
	v, exists := m.GetIf(key)

	if exists {
		if v, ok := v.(map[string]interface{}); ok {
			return v
		}
	}

	return map[string]interface{}{}
}

func (m *StringMap) M(key string) *StringMap {
	_map := m.GetOrNil(key)

	if _map == nil {
		m.Set(key, NewStringMap())

		return m.M(key)
	}

	switch _map := _map.(type) {
	case map[string]interface{}:
		m.Set(key, NewStringMapFrom(_map))
		return m.M(key)
	case *StringMap:
		return _map
	}

	return nil
}

func (m *StringMap) A(key string) *Array {
	_map := m.GetOrNil(key)

	if _map == nil {
		m.Set(key, NewArray())

		return m.A(key)
	}

	switch _map := _map.(type) {
	case []interface{}:
		m.Set(key, NewArray().Add(_map...))
		return m.A(key)
	case *Array:
		return _map
	}

	return nil
}

func (m StringMap) MarshalJSON() ([]byte, error) {

	return json.Marshal(map[string]interface{}(m))
}

func (m StringMap) JSON() string {
	if b, err := json.Marshal(map[string]interface{}(m)); err == nil {
		return string(b)
	}
	return "{}"
}

func (m StringMap) ToMap() map[string]interface{} {
	return map[string]interface{}(m)
}
