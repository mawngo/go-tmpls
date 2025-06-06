package internal

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"time"
)

func NewBuiltinFuncMap(excludes ...string) template.FuncMap {
	builtin := template.FuncMap{
		// Pack value pairs into a map.
		// Example: dict "Users" .MostPopular "Current" .CurrentUser.
		"dict":   dict,
		"dig":    dig,
		"strval": strval,
		"N":      N,
		"int":    toInt,
		"add": func(a any, i ...any) int {
			sum := toInt(a)
			for _, b := range i {
				sum += toInt(b)
			}
			return sum
		},
		"sub": func(a, b any) int { return toInt(a) - toInt(b) },
		"div": func(a, b any) int {
			return toInt(a) / toInt(b)
		},
		"mul": func(a any, i ...any) int {
			total := toInt(a)
			for _, b := range i {
				total *= toInt(b)
			}
			return total
		},
		"min": func(a any, i ...any) int {
			m := toInt(a)
			for _, b := range i {
				m = min(m, toInt(b))
			}
			return m
		},
		"max": func(a any, i ...any) int {
			m := toInt(a)
			for _, b := range i {
				m = max(m, toInt(b))
			}
			return m
		},
		"upper":    strings.ToUpper,
		"lower":    strings.ToLower,
		"title":    strings.ToTitle,
		"date":     date,
		"datetime": datetime,
		"ternary":  ternary,
	}
	for _, name := range excludes {
		delete(builtin, name)
	}
	return builtin
}

// N create a pseudo slice for range over number in template.
// The second parameter v can be used for preserving the dot type.
func N(n int, v ...any) []any {
	arr := make([]any, 0, n)
	if len(v) == 0 {
		for i := 0; i < n; i++ {
			arr = append(arr, i)
		}
		return arr
	}
	for i := 0; i < n; i++ {
		arr = append(arr, v[0])
	}
	return arr
}

// strval convert value to string.
func strval(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// date format time, default format is time.DateOnly.
func date(v time.Time, format ...string) string {
	if len(format) == 0 {
		return v.Format(time.DateOnly)
	}
	return v.Format(format[0])
}

// datetime format time, default format is time.DateTime.
func datetime(v time.Time, format ...string) string {
	if len(format) == 0 {
		return v.Format(time.DateTime)
	}
	return v.Format(format[0])
}

// dict https://github.com/Masterminds/sprig.
// Creating dictionaries is done by calling the dict function and passing it a list of pairs.
func dict(v ...any) map[string]any {
	dict := map[string]any{}
	lenv := len(v)
	for i := 0; i < lenv; i += 2 {
		key := strval(v[i])
		if i+1 >= lenv {
			dict[key] = ""
			continue
		}
		dict[key] = v[i+1]
	}
	return dict
}

// dig traverses a nested set of dicts, selecting keys from a list of values.
// It returns a default value if any of the keys are not found at the associated dict.
func dig(ps ...any) (any, error) {
	if len(ps) < 3 {
		panic("dig needs at least three arguments")
	}
	dict := ps[len(ps)-1].(map[string]any)
	def := ps[len(ps)-2]
	ks := make([]string, len(ps)-2)
	for i := 0; i < len(ks); i++ {
		ks[i] = ps[i].(string)
	}

	return digFromDict(dict, def, ks)
}

func digFromDict(dict map[string]any, d any, ks []string) (any, error) {
	k, ns := ks[0], ks[1:]
	step, has := dict[k]
	if !has {
		return d, nil
	}
	if len(ns) == 0 {
		return step, nil
	}
	return digFromDict(step.(map[string]any), d, ns)
}

// ternary returns the first value if the last value is true, otherwise returns the second value.
// true | ternary "b" "c" => "b".
// false | ternary "b" "c" => "c".
func ternary(vt any, vf any, v bool) any {
	if v {
		return vt
	}

	return vf
}

// From html/template/content.go
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

func toInt(i any) int {
	i = indirect(i)
	switch s := i.(type) {
	case int:
		return s
	case time.Weekday:
		return int(s)
	case time.Month:
		return int(s)
	case int64:
		return int(s)
	case int32:
		return int(s)
	case int16:
		return int(s)
	case int8:
		return int(s)
	case uint:
		return int(s)
	case uint64:
		return int(s)
	case uint32:
		return int(s)
	case uint16:
		return int(s)
	case uint8:
		return int(s)
	case float64:
		return int(s)
	case float32:
		return int(s)
	case nil:
		return 0
	default:
		panic(fmt.Errorf("unable to cast %#v of type %T to int", i, i))
	}
}
