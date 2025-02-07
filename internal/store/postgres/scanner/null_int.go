package scanner

import "github.com/jackc/pgtype"

func ScanInt64(value *int64) any {
	return ScanFunc(func(src any) error {
		if src == nil {
			return nil
		}
		t := pgtype.Int8{}
		err := t.Scan(src)
		if err != nil {
			return err
		}
		*value = t.Int
		return nil
	})
}


// ScanInt safely scans an SQL value into an *int.
func ScanInt(value *int) any {
	return ScanFunc(func(src any) error {
		if src == nil {
			return nil
		}
		t := pgtype.Int4{} // PostgreSQL int is typically Int4
		err := t.Scan(src)
		if err != nil {
			return err
		}
		*value = int(t.Int)
		return nil
	})
}