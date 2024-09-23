package page

import (
	"math"
	"net/url"
	"strconv"
)

// D Convenient shorthand for map[string]any
type D map[string]any

type Pageable interface {
	_dontImplThisInterface()
	HasNext() bool
	HasPrevious() bool
	NextPage() int
	PreviousPage() int
	CurrentPage() int
	ElementsPerPage() int
}

type Page[T any] struct {
	Data          []T
	TotalElements int64
	TotalPages    int
	Size          int
	Page          int
	Search        string
	Sorts         Sorts
	URL           *url.URL
}

func NewPage[T any](data []T, total int64, p Paginator) Page[T] {
	return Page[T]{
		Data:          data,
		TotalElements: total,
		TotalPages:    int(math.Ceil(float64(total) / float64(p.Size))),
		Page:          p.Page,
		Search:        p.Search,
		Size:          p.Size,
		Sorts:         p.Sorts,
		URL:           p.URL,
	}
}

func (p Page[T]) _dontImplThisInterface() {
}

func (p Page[T]) HasNext() bool {
	return p.Page < p.TotalPages
}

func (p Page[T]) CurrentPage() int {
	return p.Page
}

func (p Page[T]) ElementsPerPage() int {
	return p.Size
}

func (p Page[T]) HasPrevious() bool {
	return p.Page > 1
}

func (p Page[T]) NextPage() int {
	if !p.HasNext() {
		return p.Page
	}
	return p.Page + 1
}

func (p Page[T]) PreviousPage() int {
	if !p.HasPrevious() {
		return p.Page
	}
	return p.Page - 1
}

func (p Page[T]) PathToNext() string {
	return p.PathToPage(p.NextPage())
}

func (p Page[T]) PathToPrevious() string {
	return p.PathToPage(p.PreviousPage())
}

func (p Page[T]) PathToFirst() string {
	return p.PathToPage(FirstPageNumber)
}

func (p Page[T]) PathToLast() string {
	return p.PathToPage(p.TotalPages)
}

// PathToPage returns the URL path for the given page number.
func (p Page[T]) PathToPage(page int) string {
	query := p.URL.Query()
	query.Del(ParamPage)
	if page > 1 {
		query.Set(ParamPage, strconv.Itoa(page))
	}
	if q := query.Encode(); q != "" {
		return p.URL.Path + "?" + q
	}
	return p.URL.Path
}

// PathToSize returns the URL path for the given size.
// Changing size will reset the page to 1.
func (p Page[T]) PathToSize(size int) string {
	query := p.URL.Query()
	query.Del(ParamPage)
	query.Del(ParamSize)
	if size > 1 {
		query.Set(ParamSize, strconv.Itoa(size))
	}
	if q := query.Encode(); q != "" {
		return p.URL.Path + "?" + q
	}
	return p.URL.Path
}

// PathToSort returns the URL path for the given sort.
// Changing sorts will reset the page to 1.
func (p Page[T]) PathToSort(sorts ...string) string {
	query := p.URL.Query()
	query.Del(ParamPage)
	query.Del(ParamSort)
	if len(sorts) > 0 {
		query[ParamSort] = sorts
	}
	if q := query.Encode(); q != "" {
		return p.URL.Path + "?" + q
	}
	return p.URL.Path
}

// PathToSearch returns the URL path for the given search.
// Changing search will reset the page to 1 and unset the sort.
func (p Page[T]) PathToSearch(search string) string {
	query := p.URL.Query()
	query.Del(ParamPage)
	query.Del(ParamSort)
	query.Del(ParamSearch)
	if search != "" {
		query.Set(ParamSearch, search)
	}
	if q := query.Encode(); q != "" {
		return p.URL.Path + "?" + q
	}
	return p.URL.Path
}

// PathWithParam returns the URL path with additional query param appended.
func (p Page[T]) PathWithParam(param string, values ...string) string {
	query := p.URL.Query()
	if _, ok := query[param]; !ok {
		query[param] = make([]string, 0, len(values))
	}
	query[param] = append(query[param], values...)
	if q := query.Encode(); q != "" {
		return p.URL.Path + "?" + q
	}
	return p.URL.Path
}

// PathWithQuery returns the URL path with query string appended.
func (p Page[T]) PathWithQuery(queryString string) string {
	query := p.URL.RawQuery
	if query == "" {
		return p.URL.Path + "?" + queryString
	}
	return p.URL.Path + "?" + query + "&" + queryString
}
