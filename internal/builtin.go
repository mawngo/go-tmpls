package internal

import (
	"fmt"
	"reflect"
	"strings"
)

func NewBuiltinFuncMap(excludes ...string) map[string]any {
	builtin := map[string]any{
		"dict":    dict,
		"dig":     dig,
		"strval":  strval,
		"until":   until,
		"ternary": ternary,

		"int": toInt,
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

		"float64": toFloat64,
		"addf": func(a any, i ...any) float64 {
			sum := toFloat64(a)
			for _, b := range i {
				sum += toFloat64(b)
			}
			return sum
		},
		"subf": func(a, b any) float64 { return toFloat64(a) - toFloat64(b) },
		"divf": func(a, b any) float64 {
			return toFloat64(a) / toFloat64(b)
		},
		"mulf": func(a any, i ...any) float64 {
			total := toFloat64(a)
			for _, b := range i {
				total *= toFloat64(b)
			}
			return total
		},
		"minf": func(a any, i ...any) float64 {
			m := toFloat64(a)
			for _, b := range i {
				m = min(m, toFloat64(b))
			}
			return m
		},
		"maxf": func(a any, i ...any) float64 {
			m := toFloat64(a)
			for _, b := range i {
				m = max(m, toFloat64(b))
			}
			return m
		},

		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.ToTitle,

		"datefmt":       datefmt,
		"datetimefmt":   datetimefmt,
		"dateInZone":    dateInZone,
		"date":          date,
		"now":           now,
		"toDate":        toDate,
		"duration":      duration,
		"durationRound": durationRound,
	}
	for _, name := range excludes {
		delete(builtin, name)
	}
	return builtin
}

// until create a pseudo slice for range over number in the template.
// The second parameter v can be used for preserving the dot type.
func until(n int, v ...any) []any {
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

// dict https://github.com/Masterminds/sprig.
// Creating dictionaries is done by calling the dict function and passing it a list of pairs.
// Example: dict "Users" .MostPopular "Current" .CurrentUser.
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
