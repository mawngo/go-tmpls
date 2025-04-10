package page

import (
	"github.com/mawngo/go-tmpls/page/simplepage"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	FirstPageNumber = 1
	DefaultPageSize = 24
)

const (
	ParamPage   = "page"
	ParamSize   = "size"
	ParamSort   = "sort"
	ParamSearch = "q"
)

// Paginator represent a page request.
type Paginator struct {
	simplepage.Paginator
	// URL request URL.
	URL *url.URL
	// query request query params.
	query url.Values
}

// Query return given query param value.
func (p Paginator) Query(name string) string {
	return p.query.Get(name)
}

// Search return value of ParamSearch query param, trimmed.
func (p Paginator) Search() string {
	return strings.TrimSpace(p.Query(ParamSearch))
}

// NewDefaultPaginator returns new paginator with default values.
func NewDefaultPaginator(url *url.URL, sorts ...string) Paginator {
	return Paginator{
		Paginator: simplepage.Paginator{
			PageNumber: FirstPageNumber,
			Size:       DefaultPageSize,
			Sorts:      simplepage.NewSorts(sorts...),
		},
		query: url.Query(),
		URL:   url,
	}
}

// NewPaginator returns a new paginator from request and optionally default sorts.
func NewPaginator(req *http.Request, sorts ...string) Paginator {
	p := NewDefaultPaginator(req.URL, sorts...)
	query := p.query
	if page := query.Get(ParamPage); page != "" {
		if pageNumber, err := strconv.Atoi(page); err == nil {
			p.PageNumber = max(pageNumber, FirstPageNumber)
		}
	}
	if size := query.Get(ParamSize); size != "" {
		if pageSize, err := strconv.Atoi(size); err == nil {
			p.Size = max(pageSize, 1)
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
		p.Sorts = simplepage.NewSorts(s...)
	}
	return p
}
