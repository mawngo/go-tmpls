package simplepage

import "math"

var _ Paged[any] = (*Page[any])(nil)
var _ Paged[any] = (*Slice[any])(nil)

// Paged is the minimal interface for Pagination.
// All pages results struct implement this interface.
//
// This interface comes handy when you need to support multiple pagination strategies
// in a single api (which you should not, and which is what I did).
type Paged[T any] interface {
	GetItems() []T
	GetSorts() []Sort
	GetPageable() Pageable
}

// Slice is the basic paginated data without total items count.
type Slice[T any] struct {
	Items   []T  `json:"items"`
	HasNext bool `json:"hasNext"`
	HasPrev bool `json:"hasPrev"`

	PageNumber int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	Sorts      []Sort `json:"sorts,omitempty"`
	IsUnpaged  bool   `json:"isUnpaged,omitempty,omitzero"`
}

func (p Slice[T]) GetItems() []T {
	return p.Items
}

func (p Slice[T]) GetSorts() []Sort {
	return p.Sorts
}

// GetPageable reconstruct [Paging] from this Slice data,
// used for further processing Items and constructing new Slice.
func (p Slice[T]) GetPageable() Pageable {
	return Paging{
		Unpaged: p.IsUnpaged,
		Page:    p.PageNumber,
		Size:    p.PageSize,
		Sorts:   p.Sorts,
	}
}

// NewSlice create new Slice.
func NewSlice[T any](pageable Pageable, items []T, hasNext bool) Slice[T] {
	if pageable.IsUnpaged() {
		pageable = Paging{
			Unpaged: true,
			Page:    DefaultPageNumber,
			Size:    max(len(items), DefaultPageSize),
			Sorts:   pageable.PageSorts(),
		}
	}
	pageNumber := pageable.PageNumber()
	return Slice[T]{
		Items:      items,
		HasNext:    hasNext,
		HasPrev:    pageNumber > 1,
		PageSize:   pageable.PageSize(),
		PageNumber: pageNumber,
		Sorts:      pageable.PageSorts(),
		IsUnpaged:  pageable.IsUnpaged(),
	}
}

// Page represents paged data.
type Page[T any] struct {
	Slice[T]
	TotalPages int   `json:"totalPages"`
	TotalItems int64 `json:"totalItems"`
}

// NewPage create new Page.
func NewPage[T any](pageable Pageable, items []T, totalItems int64) Page[T] {
	if pageable.IsUnpaged() {
		pageable = Paging{
			Unpaged: true,
			Page:    DefaultPageNumber,
			Size:    max(len(items), DefaultPageSize),
			Sorts:   pageable.PageSorts(),
		}
		totalItems = max(totalItems, int64(len(items)))
	}

	pageSize := pageable.PageSize()
	pageNumber := pageable.PageNumber()
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	return Page[T]{
		Slice: Slice[T]{
			Items:      items,
			HasNext:    pageNumber < totalPages,
			HasPrev:    pageNumber > 1,
			PageSize:   pageSize,
			PageNumber: pageNumber,
			Sorts:      pageable.PageSorts(),
			IsUnpaged:  pageable.IsUnpaged(),
		},
		TotalPages: totalPages,
		TotalItems: totalItems,
	}
}
