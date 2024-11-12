package scanner

import "github.com/jackc/pgtype"

func ScanTimestamp(ref *int64) any {
	return ScanFunc(func(src any) error {
		t := pgtype.Timestamptz{}
		err := t.Scan(src)
		if err != nil {
			return err
		}
		v := t.Time.UnixMilli()
		*ref = v
		return nil
	})
}
