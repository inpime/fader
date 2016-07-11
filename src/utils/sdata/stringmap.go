package sdata

type StringMap map[string]interface{}

func NewStringMap() *StringMap {
	m := StringMap(make(map[string]interface{}))
	return &m
}

func NewStringMapFrom(v interface{}) (m *StringMap) {
	switch v := v.(type) {
	case map[string]interface{}:
		_m := StringMap(v)
		return &_m
	case *StringMap:
		m = v
	}

	return
}

func (m *StringMap) Set(key string, value interface{}) *StringMap {
	(*m)[key] = value

	return m
}

// Get alias GetOrNil
func (m StringMap) Get(key string) (value interface{}) {

	return m.GetOrNil(key)
}

func (m StringMap) GetIf(key string) (value interface{}, exists bool) {
	value, exists = m[key]
	return
}

func (m StringMap) GetOrNil(key string) interface{} {
	if value, exists := m.GetIf(key); exists {
		return value
	}

	return nil
}

func (m *StringMap) Remove(key string) *StringMap {
	delete(*m, key)
	return m
}

func (m StringMap) Size() int {
	return len(m)
}

func (m StringMap) Keys() []string {
	keys := make([]string, m.Size())
	count := 0
	for key := range m {
		keys[count] = key
		count++
	}
	return keys
}

func (m StringMap) Values() []interface{} {
	values := make([]interface{}, m.Size())
	count := 0
	for _, value := range m {
		values[count] = value
		count++
	}
	return values
}
