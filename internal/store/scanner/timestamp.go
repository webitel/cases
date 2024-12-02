package scanner

import (
	"time"
)

func ScanTimestamp(ref *int64) any {
	return ScanFunc(func(src any) error {
		if src == nil {
			return nil
		}
		var res int64
		switch val := src.(type) {
		case []byte:
			t, err := time.Parse("2006-01-02 15:04:05.999999", string(val))
			if err != nil {
				return err
			}
			res = t.UnixMilli()
		case string:
			t, err := time.Parse("2006-01-02 15:04:05.999999", val)
			if err != nil {
				return err
			}
			res = t.UnixMilli()
		case int64:
			res = val
		}
		*ref = res
		return nil
	})
}
