package internal

import (
	"encoding/json"
	"math"
	"math/rand"
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

// randInt returns a random integer between min (inclusive) and max (exclusive).
func randInt(mini, maxi int) int { return rand.Intn(maxi-mini) + mini }

// round returns a float value with the remainder rounded to the given number to digits after the decimal point.
func round(a any, p int, rOpt ...float64) float64 {
	roundOn := .5
	if len(rOpt) > 0 {
		roundOn = rOpt[0]
	}
	val := toFloat64(a)
	places := toFloat64(p)

	var round float64
	pow := math.Pow(10, places)
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return round / pow
}
