package sdata

import ()

type Map map[interface{}]interface{}

func NewMap() *Map {
	m := Map(make(map[interface{}]interface{}))
	return &m
}

func NewMapFrom(v interface{}) (m *Map) {
	switch v := v.(type) {
	case map[interface{}]interface{}:
		_m := Map(v)
		return &_m
	case *Map:
		m = v
	}

	return
}

func (m *Map) Set(key, value interface{}) *Map {
	(*m)[key] = value

	return m
}

func (m Map) Get(key interface{}) (value interface{}) {

	return m.GetOrNil(key)
}

func (m Map) GetIf(key interface{}) (value interface{}, exists bool) {
	value, exists = m[key]
	return
}

func (m Map) GetOrNil(key interface{}) interface{} {
	if value, exists := m.GetIf(key); exists {
		return value
	}

	return nil
}

func (m *Map) Remove(key interface{}) *Map {
	delete(*m, key)
	return m
}

func (m Map) Size() int {
	return len(m)
}

func (m Map) Keys() []interface{} {
	keys := make([]interface{}, m.Size())
	count := 0
	for key := range m {
		keys[count] = key
		count++
	}
	return keys
}

func (m Map) Values() []interface{} {
	values := make([]interface{}, m.Size())
	count := 0
	for _, value := range m {
		values[count] = value
		count++
	}
	return values
}
