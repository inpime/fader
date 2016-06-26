package utils

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"strconv"
)

type Stringer interface {
	String() string
}

// M Map map storage
type M map[string]interface{}

// TODO: add function Valid (see https://github.com/go-playground/validator)

// Map map storage
func Map(m ...interface{}) M {
	if len(m) == 0 {
		return M(map[string]interface{}{})
	}

	if len(m) == 1 {

		return M(map[string]interface{}{}).LoadFrom(m[0])
	}

	_m := Map()

	for _, v := range m {
		_m.LoadFrom(v)
	}

	return _m
}

func (m M) LoadFrom(v interface{}) M {
	switch v := v.(type) {
	case string:
		if len(v) > 0 {
			err := json.Unmarshal([]byte(v), &m)
			if err != nil {
				logrus.WithError(err).Error("map: load from json")
			}
		}
	case map[string]interface{}:
		m.LoadFromMapStrIfac(v)
	case map[interface{}]interface{}:
		m.LoadFromMapIfacIfac(v)
	case M:
		m.LoadFromM(v)
	default:
		// not supported
		panic("not supported")
	}

	return m
}

func (m M) LoadFromM(v M) {
	for _k, _v := range v {
		m.Set(_k, _v)
	}
}

func (m M) LoadFromMapStrIfac(v map[string]interface{}) {
	for _k, _v := range v {
		m.Set(_k, _v)
	}
}

func (m M) LoadFromMapIfacIfac(v map[interface{}]interface{}) {
	for _k, _v := range v {
		m.Set(toString(_k), _v)
	}
}

//

func (m M) Strings(k string) []string {
	v := m.GetOrNil(k)

	if v, ok := v.([]string); ok {
		return v
	}

	if v, ok := v.([]interface{}); ok {
		m.Set(k, NewA(v).(AStrings).Array())

		return m.Strings(k)
	}

	m.Set(k, []string{})

	return m.Strings(k)
}

func (m M) Map(k string) map[string]interface{} {
	v := m.GetOrNil(k)

	if v, ok := v.(map[string]interface{}); ok {
		return v
	}

	if v, ok := v.(map[interface{}]interface{}); ok {

		m.Set(k, Map(v))
		return m.Map(k)
	}

	m.Set(k, map[string]interface{}{})

	return m.Map(k)
}

func (m M) Keys(k string) (keys []string) {

	for key, _ := range m.Map(k) {
		keys = append(keys, key)
	}

	return
}

func (m M) M(k string) M {
	v := m.GetOrNil(k)

	if v, ok := v.(map[string]interface{}); ok {
		return M(v)
	}

	m.Set(k, map[string]interface{}{})

	return m.M(k)
}

// Set set value by key
func (m M) Set(k string, v interface{}) M {
	m[k] = v

	return m
}

func (m *M) SetPtr(k string, v interface{}) *M {
	(*m)[k] = v

	return m
}

// Get get value by key
func (m M) Get(k string) interface{} {

	return m[k]
}

func (m M) Delete(k string) {
	delete(m, k)
}

// GetOrNil get value by key
// If not exist retern nil
func (m M) GetOrNil(k string) interface{} {
	if v, exists := m[k]; exists {
		return v
	}

	return nil
}

// Include if exist value return true
func (m M) Include(k string) bool {
	v := m.GetOrNil(k)

	if v == nil {
		return false
	}

	return true
}

func (m M) Bool(k string) bool {
	v := m.GetOrNil(k)

	if v == nil {
		return false
	}

	switch t := v.(type) {
	case string:
		// for binding html form checkboxes
		if t == "on" {
			return true
		}

		b, err := strconv.ParseBool(t)
		if err != nil {
			return false
		}

		return b
	case bool:
		return t
		// TODO: from integers
	}

	return false
}

// String return value as string
func (m M) String(k string) string {
	v := m.GetOrNil(k)

	if v == nil {
		return ""
	}

	return toString(v)
}

// Int64 return value as int64
func (m M) Int64(k string) int64 {
	v := m.GetOrNil(k)

	if v == nil {
		return 0
	}

	switch t := v.(type) {
	case int:
		return int64(t)
	case int16:
		return int64(t)
	case int32:
		return int64(t)
	case int64:
		return t
	case uint:
		return int64(t)
	case uint8:
		return int64(t)
	case uint16:
		return int64(t)
	case uint32:
		return int64(t)
	case uint64:
		return int64(t)
	case float64:
		return int64(t)
	case float32:
		return int64(t)
	case string:

		if i, err := strconv.ParseInt(t, 10, 64); err == nil {
			return i
		}
	}

	return 0
}

// Int return value as Int
func (m M) Int(k string) int {
	v := m.GetOrNil(k)

	if v == nil {
		return 0
	}

	switch t := v.(type) {
	case int:
		return t
	case int16:
		return int(t)
	case int32:
		return int(t)
	case int64:
		return int(t)
	case uint:
		return int(t)
	case uint8:
		return int(t)
	case uint16:
		return int(t)
	case uint32:
		return int(t)
	case uint64:
		return int(t)
	case float64:
		return int(t)
	case float32:
		return int(t)
	case string:

		if i, err := strconv.Atoi(t); err == nil {
			return i
		}
	}

	return 0
}

func (m M) Float(k string) float64 {
	return m.Float64(k)
}

// Float64 return value as float64
func (m M) Float64(k string) float64 {
	v := m.GetOrNil(k)

	if v == nil {
		return 0
	}

	switch t := v.(type) {
	case int:
		return float64(t)
	case int16:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case uint:
		return float64(t)
	case uint8:
		return float64(t)
	case uint16:
		return float64(t)
	case uint32:
		return float64(t)
	case uint64:
		return float64(t)
	case float64:
		return float64(t)
	case float32:
		return float64(t)
	case string:

		if i, err := strconv.ParseFloat(t, 64); err == nil {
			return i
		}
	}

	return 0
}

func (m M) ToMap() map[string]interface{} {
	return map[string]interface{}(m)
}
