package scanner

import "github.com/jackc/pgtype"

func ScanText(value *string) any {
	return ScanFunc(func(src any) error {
		t := pgtype.Text{}
		err := t.Scan(src)
		if err != nil {
			return err
		}
		*value = t.String
		return nil
	})
}
