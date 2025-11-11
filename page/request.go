package page

import (
	"github.com/mawngo/go-tmpls/v2/simplepage"
	"net/url"
	"strconv"
	"strings"
)

const (
	ParamPage   = "page"
	ParamSize   = "size"
	ParamSort   = "sorts"
	ParamSearch = "q"
)

var _ Pageable = (*Paging)(nil)

// Pageable interface for requesting/constructing a page.
type Pageable interface {
	simplepage.Pageable

	Query(name string) string
	Search() string

	URL() *url.URL
	Queries() url.Values
}

// Paging represent a page request.
type Paging struct {
	simplepage.Paging
	// URL request URL.
	url *url.URL
	// queries request query params.
	queries url.Values
}

// Query return given query param value.
func (p Paging) Query(name string) string {
	return p.queries.Get(name)
}

// Search return value of ParamSearch query param, trimmed.
func (p Paging) Search() string {
	return strings.TrimSpace(p.queries.Get(ParamSearch))
}

func (p Paging) Queries() url.Values {
	return p.queries
}

func (p Paging) URL() *url.URL {
	return p.url
}

// NewDefaultPaging returns [Paging] with default values.
func NewDefaultPaging(url *url.URL, sorts ...string) Paging {
	return Paging{
		Paging: simplepage.NewPaging(
			simplepage.DefaultPageNumber,
			simplepage.DefaultPageSize,
			simplepage.NewSorts(sorts)...,
		),
		queries: url.Query(),
		url:     url,
	}
}

// NewPaging returns a new paginator from the request and optionally default sorts.
func NewPaging(url *url.URL, sorts ...string) Paging {
	p := NewDefaultPaging(url, sorts...)
	query := p.queries
	if page := query.Get(ParamPage); page != "" {
		if pageNumber, err := strconv.Atoi(page); err == nil {
			//nolint:staticcheck
			p.Page = pageNumber
		}
	}
	if size := query.Get(ParamSize); size != "" {
		if pageSize, err := strconv.Atoi(size); err == nil {
			//nolint:staticcheck
			p.Size = pageSize
		}
	}
	if !query.Has(ParamSort) {
		return p
	}

	s := make([]string, 0, len(query[ParamSort]))
	for _, sorts := range query[ParamSort] {
		for _, sort := range strings.Split(sorts, ",") {
			sort = strings.TrimSpace(sort)
			if len(sort) > 0 {
				s = append(s, sort)
			}
		}
	}
	if len(s) > 0 {
		//nolint:staticcheck
		p.Sorts = simplepage.NewSorts(s)
	}
	return p
}
