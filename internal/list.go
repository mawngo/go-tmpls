package internal

import (
	"fmt"
	"reflect"
)

func list(v ...any) []any {
	return v
}

// concat Concatenate arbitrary number of lists into one.
//
// Example usage: concat $myList (list 6 7) (list 8)
func concat(lists ...any) any {
	var res []any
	for _, list := range lists {
		tp := reflect.TypeOf(list).Kind()
		switch tp {
		case reflect.Slice, reflect.Array:
			l2 := reflect.ValueOf(list)
			for i := 0; i < l2.Len(); i++ {
				res = append(res, l2.Index(i).Interface())
			}
		default:
			panic(fmt.Sprintf("Cannot concat type %s as list", tp))
		}
	}
	return res
}
