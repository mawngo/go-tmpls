package page

import "strings"

func NewSorts(sorts ...string) Sorts {
	if len(sorts) == 0 {
		return nil
	}
	parseds := make([]string, 0, len(sorts))
	for _, sort := range sorts {
		if strings.HasSuffix(sort, " desc") || strings.HasSuffix(sort, " asc") {
			parseds = append(parseds, sort)
			continue
		}
		if sort[0] == '-' {
			parseds = append(parseds, sort[1:]+" desc")
			continue
		}
		parseds = append(parseds, sort+" asc")
	}
	return parseds
}

type Sorts []string

func (s Sorts) String() string {
	return strings.Join(s, ", ")
}

func (s Sorts) Label() string {
	if len(s) == 0 {
		return ""
	}
	str := strings.Join(s, ", ")
	str = strings.ReplaceAll(str, "_", " ")
	str = strings.ReplaceAll(str, " desc", " ↓")
	str = strings.ReplaceAll(str, " asc", " ↑")
	return str
}
