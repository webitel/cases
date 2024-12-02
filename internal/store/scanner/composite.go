package scanner

import (
	"github.com/jackc/pgtype"
	"strconv"
)

func GetCompositeTextScanFunction[T any](subScanPlan []func(*T) any, into *[]*T) TextDecoder {
	return func(src []byte) error {

		rows, err := pgtype.ParseUntypedTextArray(string(src))
		if err != nil {
			return err
		}

		var (
			raw *pgtype.CompositeTextScanner
		)
		for _, row := range rows.Elements {
			var node T
			raw = pgtype.NewCompositeTextScanner(nil, []byte(row))
			for _, scan := range subScanPlan {
				scanNode := scan(&node)
				switch scanNode.(type) {
				case *string, *[]byte, TextDecoder:
					raw.ScanValue(scanNode)
				default:
					scanFunc := TextDecoder(func(src []byte) error {
						switch val := scanNode.(type) {
						case ScanFunc:
							err = val.Scan(src)
							if err != nil {
								return err
							}
						case *int64:
							resultingInt, err := strconv.ParseInt(string(src), 10, 64)
							if err != nil {
								return err
							}
							*val = resultingInt
						case *int32:
							resultingInt, err := strconv.ParseInt(string(src), 10, 32)
							if err != nil {
								return err
							}
							*val = int32(resultingInt)
						}
						return nil
					})
					raw.ScanValue(&scanFunc)
				}
			}
			scanErr := raw.Err()
			if scanErr != nil {
				return scanErr
			}
			*into = append(*into, &node)
		}
		return nil
	}
}
