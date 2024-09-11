package internal

import (
	"fmt"
	"text/template"
)

func NewBuiltinFuncMap(excludes ...string) template.FuncMap {
	builtin := template.FuncMap{
		// Pack value pairs into a map.
		// Example: dict "Users" .MostPopular "Current" .CurrentUser.
		"dict":   dict,
		"dig":    dig,
		"strval": strval,
		"N": func(n int, v ...any) []any {
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
		},
		"add": func(i ...int) int {
			a := 0
			for _, b := range i {
				a += b
			}
			return a
		},
		"sub": func(a, b int) int { return a - b },
	}
	for _, name := range excludes {
		delete(builtin, name)
	}
	return builtin
}

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

// dict https://github.com/Masterminds/sprig
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
