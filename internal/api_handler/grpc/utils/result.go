package utils

// Lister defines an interface for types that support pagination by size.
type Lister interface {
	GetSize() int
}

// GetListResult determines if there is a next page and returns the appropriate slice of items.
// If the number of items is one more than the requested size, it indicates there is a next page.
func GetListResult[C any](s Lister, items []C) (bool, []C) {
	if len(items)-1 == s.GetSize() {
		return true, items[0 : len(items)-1]
	}
	return false, items
}

// ConvertToOutputBulk converts a slice of input items to a slice of output items using the provided convertFunc.
// If any conversion fails, it returns an error.
func ConvertToOutputBulk[C any, K any](items []C, convertFunc func(C) (K, error)) ([]K, error) {
	result := make([]K, 0, len(items))
	for _, item := range items {
		out, err := convertFunc(item)
		if err != nil {
			return nil, err
		}
		result = append(result, out)
	}
	return result, nil
}

// ResolvePaging returns a slice of items up to the specified size and a boolean indicating if there is a next page.
// If size is zero or negative, all items are returned and next is false.
func ResolvePaging[C any](size int, items []C) (result []C, next bool) {
	if size > 0 && len(items) >= size {
		return items[0:size], true
	}
	return items, false
}
