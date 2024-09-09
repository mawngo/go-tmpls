package page

import (
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

type Paginator struct {
	Size   int
	Page   int
	Search string
	Sorts  Sorts
	URL    *url.URL
}

func NewDefaultPaginator(url *url.URL, sorts ...string) Paginator {
	return Paginator{
		Page:  FirstPageNumber,
		Size:  DefaultPageSize,
		Sorts: NewSorts(sorts...),
		URL:   url,
	}
}

func NewPaginator(req *http.Request, sorts ...string) Paginator {
	p := NewDefaultPaginator(req.URL, sorts...)
	query := req.URL.Query()
	if page := query.Get(ParamPage); page != "" {
		if pageNumber, err := strconv.Atoi(page); err == nil {
			p.Page = max(pageNumber, FirstPageNumber)
		}
	}
	if size := query.Get(ParamSize); size != "" {
		if pageSize, err := strconv.Atoi(size); err == nil {
			p.Size = max(pageSize, 1)
		}
	}
	if search := query.Get(ParamSearch); search != "" {
		p.Search = search
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
	p.Sorts = NewSorts(s...)
	return p
}
