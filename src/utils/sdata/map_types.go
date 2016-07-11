package sdata

import (
	"encoding/json"
)

func (m Map) String(key interface{}) string {
	return toString(m.GetOrNil(key))
}

func (m Map) Float(key interface{}) float64 {
	return toFloat64(m.GetOrNil(key))
}

func (m Map) Int(key interface{}) int {
	return toInt(m.GetOrNil(key))
}

func (m Map) Int64(key interface{}) int64 {
	return toInt64(m.GetOrNil(key))
}

func (m Map) Bool(key interface{}) bool {
	return toBool(m.GetOrNil(key))
}

func (m *Map) Map(key interface{}) map[interface{}]interface{} {
	v, exists := m.Get(key)

	if exists {
		if v, ok := v.(map[interface{}]interface{}); ok {
			return v
		}
	}

	return map[interface{}]interface{}{}
}

func (m *Map) M(key interface{}) *Map {
	_map := m.GetOrNil(key)

	if _map == nil {
		m.Set(key, NewMap())

		return m.M(key)
	}

	switch _map := _map.(type) {
	case map[interface{}]interface{}:
		m.Set(key, NewMapFrom(_map))
		return m.M(key)
	case *Map:
		return _map
	}

	return nil
}

func (m *Map) A(key interface{}) *Array {
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

func (m Map) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{}, m.Size())

	for _, key := range m.Keys() {
		data[toString(key)] = m.GetOrNil(key)
	}

	return json.Marshal(data)
}
