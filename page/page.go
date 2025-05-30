package page

import (
	"github.com/mawngo/go-tmpls/page/simplepage"
	"math"
	"net/url"
	"strconv"
	"strings"
	"unicode"
)

// D Convenient shorthand for map[string]any.
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
	simplepage.Page[T]
	// URL request URL.
	URL *url.URL
	// query request query params.
	query url.Values
}

// NewPage returns a new page from data, total items and paginator.
func NewPage[T any](data []T, total int64, p Paginator) Page[T] {
	return Page[T]{
		Page: simplepage.Page[T]{
			Data:          data,
			TotalElements: total,
			TotalPages:    max(int(math.Ceil(float64(total)/float64(p.Size))), 1),
			PageNumber:    p.PageNumber,
			Size:          p.Size,
			Sorts:         p.Sorts,
		},
		URL:   p.URL,
		query: p.query,
	}
}

func (p Page[T]) _dontImplThisInterface() {
}

// HasNext return whether this is the last page.
func (p Page[T]) HasNext() bool {
	return p.PageNumber < p.TotalPages
}

// CurrentPage return the current page number.
func (p Page[T]) CurrentPage() int {
	return p.PageNumber
}

// ElementsPerPage return the page size.
func (p Page[T]) ElementsPerPage() int {
	return p.Size
}

// HasPrevious return whether this is the first page.
func (p Page[T]) HasPrevious() bool {
	return p.PageNumber > 1
}

// NextPage return the next page number, or current page if it is the last page.
func (p Page[T]) NextPage() int {
	if !p.HasNext() {
		return p.PageNumber
	}
	return p.PageNumber + 1
}

// PreviousPage return the next page number, or current page if it is the first page.
func (p Page[T]) PreviousPage() int {
	if !p.HasPrevious() {
		return p.PageNumber
	}
	return p.PageNumber - 1
}

// PathToNext return the url to the next page.
func (p Page[T]) PathToNext() string {
	return p.PathToPage(p.NextPage())
}

// PathToPrevious return the url to the previous page.
func (p Page[T]) PathToPrevious() string {
	return p.PathToPage(p.PreviousPage())
}

// PathToFirst return the url to the first page.
func (p Page[T]) PathToFirst() string {
	return p.PathToPage(FirstPageNumber)
}

// PathToLast return the url to the last page.
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

// PathToQueryParam returns the URL path for to the given query param.
// Changing to query param will reset the page to 1 and unset the sort.
// The query param will be replaced.
func (p Page[T]) PathToQueryParam(param string, values ...string) string {
	query := p.URL.Query()
	query.Del(ParamPage)
	query.Del(ParamSort)
	query.Del(ParamSearch)
	query[param] = values
	if q := query.Encode(); q != "" {
		return p.URL.Path + "?" + q
	}
	return p.URL.Path
}

// PathWithQueryParam returns the URL path with additional query param appended.
func (p Page[T]) PathWithQueryParam(param string, values ...string) string {
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

// Query return given query param value.
func (p Page[T]) Query(name string) string {
	return p.query.Get(name)
}

// Search return value of ParamSearch query param, trimmed.
func (p Page[T]) Search() string {
	return strings.TrimSpace(p.query.Get(ParamSearch))
}

// QSearch return given query param value or its value inside searching under format <param>:<value>.
func (p Page[T]) QSearch(name string) string {
	return p.QParam(name, ParamSearch)
}

// QParam return given query param value or its value inside another param under format <param>:<value>.
func (p Page[T]) QParam(name string, searchParam string) string {
	if q := p.Query(name); q != "" {
		return q
	}
	search := p.query.Get(searchParam)
	if len(search) <= len(name)+1 {
		return ""
	}
	param := name + ":"
	index := strings.Index(search, param)
	if index < 0 {
		return ""
	}
	index += len(param)
	end := len(search)
	for i := index; i < len(search); i++ {
		r := search[i]
		if unicode.IsSpace(rune(r)) {
			end = i
			break
		}
	}
	return search[index:end]
}

// Offset return the page item offset, useful for building the database query.
// For Limit use Size.
func (p Page[T]) Offset() string {
	return strings.TrimSpace(p.Query(ParamSearch))
}
