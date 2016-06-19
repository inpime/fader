package utils

func (a strarr) Len() int {
	return len(a)
}

func (a strarr) Index(value string) int {
	if a.Len() == 0 {
		return -1
	}

	for _index, _value := range a {
		if _value == value {
			return _index
		}
	}

	return -1
}

func (a strarr) Include(value string) bool {
	return a.Index(value) != -1
}

func (a *strarr) Add(value string) A {
	if a.Include(value) {
		return a
	}

	(*a) = append((*a), value)

	return a
}

func (a *strarr) Delete(value string) A {
	i := a.Index(value)

	if i == -1 {
		return a
	}

	(*a) = (*a)[:i+copy((*a)[i:], (*a)[i+1:])]

	return a
}

func (a strarr) Array() []string {
	return a
}
