package scanner

import "github.com/jackc/pgtype"

func ScanBool(value *bool) any {
	return ScanFunc(func(src any) error {
		t := pgtype.Bool{}
		err := t.Scan(src)
		if err != nil {
			return err
		}
		*value = t.Bool
		return nil
	})
}
