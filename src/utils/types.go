package utils

import (
	"strconv"
)

func toString(v interface{}) string {
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
