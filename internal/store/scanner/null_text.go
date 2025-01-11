package scanner

import "github.com/jackc/pgtype"

func ScanText(value *string) any {
	return ScanFunc(func(src any) error {
		if src == nil {
			return nil
		}
		t := pgtype.Text{}
		err := t.Scan(src)
		if err != nil {
			return err
		}
		if t.Status == pgtype.Present {
			*value = t.String
		}
		return nil
	})
}
