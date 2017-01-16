package test

import "sync"

type testcontext struct {
	vals   []interface{}
	keys   []string
	locker sync.RWMutex
}

func (t *testcontext) index(_key string) int {
	for index, key := range t.keys {
		if _key == key {
			return index
		}
	}

	return -1
}

func (t testcontext) Get(key string) interface{} {
	i := t.index(key)
	if i == -1 {
		return nil
	}

	return t.vals[i]
}

func (t *testcontext) Set(key string, v interface{}) {
	t.locker.Lock()
	defer t.locker.Unlock()

	i := t.index(key)
	t.set(i, key, v)
}

func (t *testcontext) set(index int, key string, v interface{}) {
	if index == -1 {
		t.keys = append(t.keys, key)
		t.vals = append(t.vals, v)
		return
	}

	t.vals[index] = v
}
