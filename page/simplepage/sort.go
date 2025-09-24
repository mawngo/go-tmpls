package simplepage

import (
	"encoding/json"
	"strings"
	"unicode"
)

// Sort is the sorting metadata.
type Sort struct {
	Field  string `json:"field,omitempty"`
	IsDesc bool   `json:"isDesc,omitempty"`
}

// FieldStrict return [Sort.Field] if it only contains
// letter, number characters, dash, underscore and space.
func (s Sort) FieldStrict() string {
	for _, c := range s.Field {
		if c == '_' || c == '-' || unicode.IsSpace(c) {
			continue
		}
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			continue
		}
		return ""
	}
	return s.Field
}

// NewSorts create a list of [Sort] from a list of string.
// String starting with negative sign '-' indicate desc sort.
// Optionally, starting with positive sign '+' indicate asc sort.
func NewSorts(raw []string) []Sort {
	if len(raw) == 0 {
		return nil
	}

	sorts := make([]Sort, 0, len(raw))
	for _, rawSort := range raw {
		isDesc := false
		if len(rawSort) > 1 {
			switch rawSort[0] {
			case '-':
				isDesc = true
				rawSort = rawSort[1:]
			case '+':
				rawSort = rawSort[1:]
			}
		}
		if rawSort == "" {
			continue
		}
		sorts = append(sorts, Sort{Field: rawSort, IsDesc: isDesc})
	}
	return sorts
}

// Sorts is a list of [Sort], with custom binding logic.
// You should use []Sort directly, except for binding.
type Sorts []Sort

// UnmarshalText support coma separated list.
func (s *Sorts) UnmarshalText(text []byte) error {
	sorts := strings.Split(string(text), ",")
	*s = NewSorts(sorts)
	return nil
}

// UnmarshalParam support coma separated list.
func (s *Sorts) UnmarshalParam(param string) error {
	sorts := strings.Split(param, ",")
	*s = NewSorts(sorts)
	return nil
}

// UnmarshalJSON support coma separated list string or array of string.
func (s *Sorts) UnmarshalJSON(b []byte) error {
	data := string(b)
	if data == "null" {
		return nil
	}

	if len(data) > 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[len(`"`) : len(data)-len(`"`)]
		return s.UnmarshalParam(data)
	}

	var sorts []string
	if err := json.Unmarshal(b, &sorts); err != nil {
		return err
	}
	*s = NewSorts(sorts)
	return nil
}

// String return SQL sort string representation.
// Deprecated: this method is prone to SQL injection and should not be used. Will be removed in the future.
func (s Sorts) String() string {
	buff := strings.Builder{}
	for i := 0; i < len(s); i++ {
		if s[i].Field == "" {
			continue
		}
		if i > 0 {
			buff.WriteString(", ")
		}

		buff.WriteString(s[i].Field)
		if s[i].IsDesc {
			buff.WriteString(" DESC")
		} else {
			buff.WriteString(" ASC")
		}
	}
	return buff.String()
}

// Label return sort (field-strict) with direction indicated by arrow.
func (s Sorts) Label() string {
	if len(s) == 0 {
		return ""
	}
	buff := strings.Builder{}
	for i := 0; i < len(s); i++ {
		field := s[i].FieldStrict()
		if field == "" {
			continue
		}
		if i > 0 {
			buff.WriteString(", ")
		}

		buff.WriteString(strings.ReplaceAll(field, "_", " "))
		if s[i].IsDesc {
			buff.WriteString(" ↓")
		} else {
			buff.WriteString(" ↑")
		}
	}
	return buff.String()
}
