package page

import (
	"github.com/mawngo/go-tmpls/v2/simplepage"
	"net/url"
	"strconv"
	"strings"
)

const (
	DefaultPageNumber = simplepage.DefaultPageNumber
	DefaultPageSize   = simplepage.DefaultPageSize
)

// D Convenient shorthand for map[string]any.
type D map[string]any

var _ Paged[any] = (*Page[any])(nil)

// PagedData interface for casting to any [Paged] type.
type PagedData interface {
	GetPageable() simplepage.Pageable
	GetSorts() []Sort
	IsEmpty() bool
}

// Paged is the minimal interface for Pagination.
// All pages results struct implement this interface.
type Paged[T any] interface {
	simplepage.Paged[T]

	HasNext() bool
	HasPrevious() bool
	NextPage() int
	PreviousPage() int
	CurrentPage() int

	PathToNext() string
	PathToPrevious() string
	PathToFirst() string
	PathToLast() string
	PathToPage(page int) string
	PathToSize(size int) string
	PathToSort(sorts ...string) string

	PathToQueryParam(param string, values ...string) string
	PathWithQueryParam(param string, values ...string) string

	Query(name string) string
	Search() string

	URL() *url.URL
	QueryValues() url.Values
}

// Page represents a page of data.
type Page[T any] struct {
	simplepage.Page[T]
	// URL request URL.
	url *url.URL
	// queries request query params.
	queries url.Values
}

// NewPage returns a new [Page] from paginator, data, and total items count.
func NewPage[T any](p Pageable, items []T, total int64) Page[T] {
	return Page[T]{
		Page:    simplepage.NewPage(p, items, total),
		url:     p.URL(),
		queries: p.QueryValues(),
	}
}

// NewEmptyPage returns a new empty [Page].
func NewEmptyPage[T any](p Pageable) Page[T] {
	return NewPage[T](p, nil, 0)
}

// HasNext return whether this is the last page.
func (p Page[T]) HasNext() bool {
	return p.PageNumber < p.TotalPages
}

// CurrentPage return the current page number.
func (p Page[T]) CurrentPage() int {
	return p.PageNumber
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
	return p.PathToPage(simplepage.DefaultPageNumber)
}

// PathToLast return the url to the last page.
func (p Page[T]) PathToLast() string {
	return p.PathToPage(p.TotalPages)
}

// PathToPage returns the URL path for the given page number.
func (p Page[T]) PathToPage(page int) string {
	query := p.url.Query()
	query.Del(ParamPage)
	if page > 1 {
		query.Set(ParamPage, strconv.Itoa(page))
	}
	if q := query.Encode(); q != "" {
		return p.url.Path + "?" + q
	}
	return p.url.Path
}

// PathToSize returns the URL path for the given size.
// Changing the size will reset the page to 1.
func (p Page[T]) PathToSize(size int) string {
	query := p.url.Query()
	query.Del(ParamPage)
	query.Del(ParamSize)
	if size > 1 {
		query.Set(ParamSize, strconv.Itoa(size))
	}
	if q := query.Encode(); q != "" {
		return p.url.Path + "?" + q
	}
	return p.url.Path
}

// PathToSort returns the URL path for the given sort.
// Changing sorts will reset the page to 1.
func (p Page[T]) PathToSort(sorts ...string) string {
	query := p.url.Query()
	query.Del(ParamPage)
	query.Del(ParamSort)
	if len(sorts) > 0 {
		query[ParamSort] = sorts
	}
	if q := query.Encode(); q != "" {
		return p.url.Path + "?" + q
	}
	return p.url.Path
}

// PathToQueryParam returns the URL path for to the given query param.
//
// Changing to the query param will reset the page to 1 and unset the sort.
// The query param will be replaced.
func (p Page[T]) PathToQueryParam(param string, values ...string) string {
	query := p.url.Query()
	query.Del(ParamPage)
	query.Del(ParamSort)
	query.Del(ParamSearch)
	query[param] = values
	if q := query.Encode(); q != "" {
		return p.url.Path + "?" + q
	}
	return p.url.Path
}

// PathWithQueryParam returns the URL path with an additional query param appended.
//
// Does not change the page or sort.
// The query param will be appended.
func (p Page[T]) PathWithQueryParam(param string, values ...string) string {
	query := p.url.Query()
	if _, ok := query[param]; !ok {
		query[param] = make([]string, 0, len(values))
	}
	query[param] = append(query[param], values...)
	if q := query.Encode(); q != "" {
		return p.url.Path + "?" + q
	}
	return p.url.Path
}

// Query return given query param value.
func (p Page[T]) Query(name string) string {
	return p.queries.Get(name)
}

// Search return value of ParamSearch query param, trimmed.
func (p Page[T]) Search() string {
	return strings.TrimSpace(p.queries.Get(ParamSearch))
}

// QueryValues return parsed request query params.
func (p Page[T]) QueryValues() url.Values {
	return p.queries
}

// URL return request URL.
func (p Page[T]) URL() *url.URL {
	return p.url
}
