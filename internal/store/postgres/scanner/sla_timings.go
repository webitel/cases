package scanner

import (
	"fmt"
	_go "github.com/webitel/cases/api/cases"
)

func ScanMillisToTimings(
	ptr *int64,
	convert func() (*_go.Timings, error),
	assign func(*_go.Timings),
) any {
	return &timingsScanner{
		ptr:     ptr,
		convert: convert,
		assign:  assign,
	}
}

type timingsScanner struct {
	ptr     *int64
	convert func() (*_go.Timings, error)
	assign  func(*_go.Timings)
}

func (s *timingsScanner) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		*s.ptr = v
	case int32:
		*s.ptr = int64(v)
	case nil:
		*s.ptr = 0
	default:
		return fmt.Errorf("unsupported src type for millis: %T", src)
	}

	if s.ptr == nil || *s.ptr <= 0 {
		return nil
	}

	t, err := s.convert()
	if err != nil {
		return err
	}
	s.assign(t)
	return nil
}
