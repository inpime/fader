package sdata

func NewArrayFrom(v interface{}) (a *Array) {

	switch v := v.(type) {
	case []interface{}:
		_a := Array(v)
		return &_a
	case []string:
		a = NewArray()
		for _, value := range v {
			a.Add(value)
		}
	case *Array:
		a = v
	default:
		a = NewArray()
	}

	return
}
