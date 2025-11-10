package internal

import (
	"encoding/json"
	"strconv"
	"time"
)

func toNumber[T int | float64](i any) T {
	i = indirect(i)
	switch s := i.(type) {
	case int:
		return T(s)
	case int8:
		return T(s)
	case int16:
		return T(s)
	case int32:
		return T(s)
	case int64:
		return T(s)
	case uint:
		return T(s)
	case uint8:
		return T(s)
	case uint16:
		return T(s)
	case uint32:
		return T(s)
	case uint64:
		return T(s)
	case float32:
		return T(s)
	case float64:
		return T(s)
	case bool:
		if s {
			return 1
		}
		return 0
	case nil:
		return 0
	case time.Weekday:
		return T(s)
	case time.Month:
		return T(s)

	case string:
		return parseNumber[T](s)
	case json.Number:
		if s == "" {
			return 0
		}
		return parseNumber[T](string(s))
	}
	return 0
}

func parseNumber[T int | float64](s string) T {
	var t T
	switch any(t).(type) {
	case float64:
		n, _ := strconv.ParseFloat(s, 64)
		return T(n)
	case int:
		v, _ := strconv.Atoi(s)
		return T(v)
	}
	return 0
}

func toFloat64(i any) float64 {
	return toNumber[float64](i)
}

func toInt(i any) int {
	return toNumber[int](i)
}
