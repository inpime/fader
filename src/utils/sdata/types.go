package sdata

import (
	"strconv"
)

type Stringer interface {
	String() string
}

func ToString(v interface{}) string {
	return toString(v)
}

func ToInt(v interface{}) int {
	return toInt(v)
}

func ToInt64(v interface{}) int64 {
	return toInt64(v)
}

func ToFloat64(v interface{}) float64 {
	return toFloat64(v)
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}

	switch t := v.(type) {
	case int:
		return strconv.Itoa(t)
	case int16:
		return strconv.FormatInt(int64(t), 10)
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(int64(t), 10)
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case uint8:
		return strconv.FormatUint(uint64(t), 10)
	case uint16:
		return strconv.FormatUint(uint64(t), 10)
	case uint32:
		return strconv.FormatUint(uint64(t), 10)
	case uint64:
		return strconv.FormatUint(uint64(t), 10)
	case float64:
		return strconv.FormatFloat(float64(t), 'E', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(t), 'E', -1, 64)
	case string:
		return t
	case Stringer:
		return t.String()
	}

	return ""
}

func toBool(v interface{}) bool {
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

func toInt64(v interface{}) int64 {
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

func toInt(v interface{}) int {
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

func toFloat64(v interface{}) float64 {
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
