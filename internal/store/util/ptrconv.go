package util

// Int64PtrOrNil takes a function returning int64 and returns a pointer if the result is > 0.
func Int64PtrOrNil(idGetter func() int64) *int64 {
	if id := idGetter(); id > 0 {
		return &id
	}
	return nil
}

// StrPtrOrNil returns a pointer to the string if it's not empty.
func StrPtrOrNil(s string) *string {
	if s != "" {
		return &s
	}
	return nil
}

// IDPtr returns a pointer to the ID if the entity is non-nil and the ID > 0.
func IDPtr(entity interface{ GetId() int64 }) *int64 {
	if entity != nil && entity.GetId() > 0 {
		return &[]int64{entity.GetId()}[0]
	}
	return nil
}

// StringPtr returns a pointer to a non-empty string.
func StringPtr(val string) *string {
	if val != "" {
		return &val
	}
	return nil
}
