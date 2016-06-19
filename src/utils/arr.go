package utils

type AStrings interface {
	Array() []string
}

type A interface {
	Len() int
	Index(string) int
	Include(string) bool
	Add(string) A
	Delete(string) A
}

type strarr []string

var _ A = (*strarr)(nil)
var _ AStrings = (*strarr)(nil)

func NewA(v interface{}) A {
	if v, ok := v.([]string); ok {
		a := strarr(v)
		return &a
	}

	if v, ok := v.([]interface{}); ok {
		_v := []string{}
		for _, _value := range v {

			_v = append(_v, toString(_value))
		}

		return NewA(_v)
	}

	return nil
}
