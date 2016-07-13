package sdata

type Array []interface{}

func NewArray() *Array {
	return new(Array)
}

func (a *Array) Add(values ...interface{}) *Array {
	(*a) = append(*a, values...)
	return a
}

func (a Array) Get(index int) interface{} {
	if index > a.Size() {
		return nil
	}

	return a[index]
}

func (a *Array) Remove(value interface{}) *Array {
	i := a.Index(value)

	if i == -1 {
		return a
	}

	(*a) = (*a)[:i+copy((*a)[i:], (*a)[i+1:])]

	return a
}

func (a Array) Index(value interface{}) int {
	if a.Size() == 0 {
		return -1
	}

	for _index, _value := range a {
		if _value == value {
			return _index
		}
	}

	return -1
}

func (a Array) Includes(values ...interface{}) bool {
	for _, value := range values {
		found := false
		for _, element := range a {
			if element == value {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (a Array) Values() []interface{} {
	newA := make([]interface{}, len(a), len(a))
	copy(newA, a[:])
	return newA
}

func (a Array) Size() int {
	return len(a)
}
