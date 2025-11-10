package simplepage

const (
	// DefaultPageNumber the page number to be used instead when [Paging.PageNumber] is zero.
	DefaultPageNumber = 1
	// DefaultPageSize the page size to be used instead when [Paging.PageSize] is zero.
	DefaultPageSize = 20
	// MaxPageSize maximum page size, hard limit and should be enforced (except when unpaged).
	MaxPageSize = 500
)

var _ Pageable = (*Paging)(nil)

// Pageable interface for requesting/constructing a page.
type Pageable interface {
	// IsUnpaged is a special flags for disabling paging behavior.
	// Api that integrates this package should properly handle this field.
	IsUnpaged() bool
	// PageNumber returning the page number.
	// Never <= 0.
	PageNumber() int
	// PageSize returning the page size.
	// Never <= 0.
	PageSize() int
	// PageOffset returning the offset of the first item in the page.
	PageOffset() int64
	// PageSorts return sort config for this page.
	PageSorts() []Sort
}

// Paging contains paging request information for embedding in DTO.
// All fields of this struct are optional for conveniently,
// so [Paging.Page] and [Paging.Size] should not be read directly
// but using Page* getter for provided default values.
//
// API that consumes this struct must also be aware of this optionality.
// However, it is recommended to consume [Pageable] instead.
type Paging struct {
	// Unpaged is a special flags for disabling paging.
	// Api that integrates this package should properly handle this field.
	Unpaged bool `json:"-" form:"-"`
	// Deprecated: write-only, for read use [Paging.PageNumber].
	Page int `json:"page" form:"page"`
	// When integrating with gin, it can be controlled by registering a "pagesize" validator.
	// For other frameworks, you may need to check by hand or write your own integration.
	// Deprecated: write-only, for read use [Paging.PageSize].
	Size int `json:"pageSize" form:"pageSize" binding:"pagesize"`
	// Sorts is a list or sort.
	// For param, it can accept a list of values separated by comma.
	// For JSON field, it can accept a string list of values separated by comma, or list of string.
	// It is recommended to validate the sort values before using them.
	// Deprecated: write-only, for read use [SortablePaging.PageSorts].
	Sorts Sorts `json:"sorts" form:"sorts"`
}

// IsUnpaged return whether paging is disabled.
// See [Paging.Unpaged].
func (p Paging) IsUnpaged() bool {
	return p.Unpaged
}

// PageNumber returning the page number.
// Never <= 0.
func (p Paging) PageNumber() int {
	return max(p.Page, DefaultPageNumber)
}

// PageSize returning the page size.
// Never <= 0.
func (p Paging) PageSize() int {
	if p.Size > 0 {
		return min(p.Size, MaxPageSize)
	}
	return DefaultPageSize
}

// PageOffset returning the offset of the first item in the page.
func (p Paging) PageOffset() int64 {
	return int64(p.PageSize()) * (int64(p.PageNumber()) - 1)
}

// PageSorts return sorts configuration of this paging.
func (p Paging) PageSorts() []Sort {
	return p.Sorts
}

// Unsorted return a [Paging] without sorting.
func (p Paging) Unsorted() Paging {
	return Paging{
		Unpaged: p.Unpaged,
		Page:    p.Page,
		Size:    p.Size,
	}
}

// NewPaging create a [Paging].
func NewPaging(page int, size int, sorts []Sort) Paging {
	return Paging{
		Page:  page,
		Size:  size,
		Sorts: sorts,
	}
}

// NewUnpaged construct an empty [Paging].
func NewUnpaged() Paging {
	return Paging{
		Unpaged: true,
	}
}
