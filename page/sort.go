package page

import "github.com/mawngo/go-tmpls/v2/page/simplepage"

type Sorts = simplepage.Sorts
type Sort = simplepage.Sort

// NewSorts create a list of [Sort] from a list of string.
// String starting with negative sign '-' indicate desc sort.
// Optionally, starting with positive sign '+' indicate asc sort.
//
// See [simplepage.NewSorts] for more details.
func NewSorts(raw ...string) Sorts {
	return simplepage.NewSorts(raw)
}
