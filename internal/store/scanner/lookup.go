package scanner

import (
	"fmt"
	"strconv"

	"github.com/jackc/pgtype"
	_go "github.com/webitel/cases/api/cases"
)

type ScanFunc func(src any) error

func (s ScanFunc) Scan(src any) error {
	return s(src)
}

type TextDecoder func(src []byte) error

func (dec TextDecoder) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	if dec != nil {
		return dec(src)
	}
	return nil
}

func (dec TextDecoder) Scan(src any) error {
	if src == nil {
		return dec.DecodeText(nil, nil)
	}

	switch data := src.(type) {
	case string:
		return dec.DecodeText(nil, []byte(data))
	case []byte:
		text := make([]byte, len(data))
		copy(text, data)
		return dec.DecodeText(nil, text)
	}

	return fmt.Errorf("text_decoder: cannot scan %T value %[1]v into %T", src, dec)
}

func ScanRowLookup(value **_go.Lookup) any {
	return TextDecoder(func(src []byte) error {
		res := *(value)
		*(value) = nil

		if len(src) == 0 {
			return nil // NULL
		}

		if res == nil {
			res = new(_go.Lookup)
		}

		var (
			ok  bool
			str pgtype.Text
			row = []pgtype.TextDecoder{
				TextDecoder(func(src []byte) error {
					err := str.DecodeText(nil, src)
					if err != nil {
						return err
					}
					id, err := strconv.ParseInt(str.String, 10, 64)
					if err != nil {
						return err
					}
					res.Id = id
					return nil
				}),
				TextDecoder(func(src []byte) error {
					err := str.DecodeText(nil, src)
					if err != nil {
						return err
					}
					res.Name = str.String
					ok = ok || (str.String != "" && str.String != "[deleted]") // && str.Status == pgtype.Present
					return nil
				}),
			}
			raw = pgtype.NewCompositeTextScanner(nil, src)
		)

		var err error
		for _, col := range row {

			raw.ScanDecoder(col)

			err = raw.Err()
			if err != nil {
				return err
			}
		}

		if ok {
			*(value) = res
		}

		return nil
	})
}
