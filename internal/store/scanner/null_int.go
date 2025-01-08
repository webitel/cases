package scanner

import "github.com/jackc/pgtype"

func ScanInt64(value *int64) any {
	return ScanFunc(func(src any) error {
		t := pgtype.Int8{}
		err := t.Scan(src)
		if err != nil {
			return err
		}
		*value = t.Int
		return nil
	})
}
