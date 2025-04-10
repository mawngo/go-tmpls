package simplepage

// Page simplest represent of a paged data response.
type Page[T any] struct {
	// Data items of page.
	Data []T
	// TotalElements total number items.
	TotalElements int64
	// TotalPages total/maximum number of page
	TotalPages int
	// Size the size of page.
	Size int
	// PageNumber the page number, start from 1.
	PageNumber int
	// Sorts the parsed sort queries.
	Sorts Sorts
}

// Paginator simplest represent of a page request.
type Paginator struct {
	// Size the size of page.
	Size int
	// PageNumber the page number, start from 1.
	PageNumber int
	// Sorts the parsed sort queries.
	Sorts Sorts
}
