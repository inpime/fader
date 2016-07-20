package sdata

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
)

func (m *StringMap) LoadFrom(v interface{}) *StringMap {
	switch v := v.(type) {
	case string:
		if len(v) > 0 {
			err := json.Unmarshal([]byte(v), m)
			if err != nil {
				logrus.WithField("_service", "utils").WithError(err).Error("map: load from json")
			}
		}
	case map[string]interface{}:
		m.LoadFromMapStrIface(v)
	case map[interface{}]interface{}:
		m.LoadFromMapIfaceIface(v)
	case *StringMap:
		m.LoadFromStringMap(v)
	}

	return m
}

func (m *StringMap) LoadFromMapStrIface(v map[string]interface{}) *StringMap {
	for key, value := range v {
		m.Set(key, value)
	}
	return m
}

func (m *StringMap) LoadFromStringMap(v *StringMap) *StringMap {
	for _, key := range v.Keys() {
		m.Set(key, v.GetOrNil(key))
	}
	return m
}

func (m *StringMap) LoadFromMapIfaceIface(v map[interface{}]interface{}) *StringMap {
	for key, value := range v {
		m.Set(toString(key), value)
	}
	return m
}

func (m *StringMap) Clear() *StringMap {
	for _, key := range m.Keys() {
		m.Remove(key)
	}
	return m
}
