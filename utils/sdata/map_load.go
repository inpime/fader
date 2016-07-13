package sdata

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
)

func (m *Map) LoadFrom(v interface{}) *Map {
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
	case *Map:
		m.LoadFromMap(v)
	}

	return m
}

func (m *Map) LoadFromMapStrIface(v map[string]interface{}) *Map {
	for key, value := range v {
		m.Set(key, value)
	}
	return m
}

func (m *Map) LoadFromMap(v *Map) *Map {
	for _, key := range v.Keys() {
		m.Set(key, v.GetOrNil(key))
	}
	return m
}

func (m *Map) LoadFromMapIfaceIface(v map[interface{}]interface{}) *Map {
	for key, value := range v {
		m.Set(toString(key), value)
	}
	return m
}
