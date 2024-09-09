package internal

import (
	"errors"
	"text/template"
)

func NewBuiltinFuncMap(excludes ...string) template.FuncMap {
	builtin := template.FuncMap{
		// Pack value pairs into a map.
		// Example: dict "Users" .MostPopular "Current" .CurrentUser.
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		// Subtract value in pipeline.
		// Example: $a | sub $b.
		"sub": func(y, x int) int {
			return x - y
		},
		// Add value in pipeline.
		// Example: $a | add $b.
		"add": func(y, x int) int {
			return x + y
		},
	}
	for _, name := range excludes {
		delete(builtin, name)
	}
	return builtin
}
