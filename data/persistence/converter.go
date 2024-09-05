package persistence

import (
	"github.com/neutrinocorp/geck/data"
	"github.com/neutrinocorp/geck/internal/converter"
)

// ConvertPage converts a page populated with type A to a page populated with type B.
func ConvertPage[A, B any](src data.Page[A], convertFunc converter.ConvertFunc[A, B]) data.Page[B] {
	return data.Page[B]{
		PreviousPageToken: src.PreviousPageToken,
		NextPageToken:     src.NextPageToken,
		TotalItems:        src.TotalItems,
		Items:             converter.ConvertMany(src.Items, convertFunc),
	}
}
