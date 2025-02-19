package scanner

import (
	"fmt"
	"time"
)

func ScanTimestamp(ref *int64) any {
	return ScanFunc(func(src any) error {
		if src == nil {
			return nil
		}
		var (
			res           int64
			tryAllLayouts = func(in string) (int64, error) {
				layouts := []string{
					time.RFC3339,
					time.RFC1123,
					time.RFC1123Z,
					time.RFC822,
					time.RFC822Z,
					time.RFC850,
					time.ANSIC,
					time.UnixDate,
					time.RubyDate,
					"2006-01-02 15:04:05.000-07",
					"2006-01-02 15:04:05",
					"02 Jan 2006 15:04:05 MST",
					"2006-01-02 15:04:05.999999",
					"2006-01-02 15:04:05.000-07",
					"2006-01-02 15:04:05.000000-07",
					"2006-01-02 15:04:05.00-07",
					"2006-01-02 15:04:05.999999-07",
				}

				// Try parsing with each layout
				var parsedTime time.Time
				var err error
				for _, layout := range layouts {
					parsedTime, err = time.Parse(layout, in)
					if err == nil {
						return parsedTime.UnixMilli(), nil
					}
				}
				return 0, fmt.Errorf("invalid date format, %s", in)
			}
		)
		switch val := src.(type) {
		case []byte:

			out, err := tryAllLayouts(string(val))
			if err != nil {
				return err
			}
			res = out
		case string:
			out, err := tryAllLayouts(val)
			if err != nil {
				return err
			}
			res = out
		case int64:
			res = val
		case time.Time:
			res = val.UnixMilli()
		}
		*ref = res
		return nil
	})
}
