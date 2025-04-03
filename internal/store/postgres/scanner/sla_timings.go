package scanner

import (
	"github.com/jackc/pgtype"
	"github.com/webitel/cases/api/cases"
	"regexp"
	"strconv"
	"strings"
)

func ScanTimingsString(out **cases.Timings) any {
	return ScanFunc(func(src any) error {
		if src == nil {
			*out = &cases.Timings{}
			return nil
		}

		t := pgtype.Text{}
		if err := t.Scan(src); err != nil {
			*out = &cases.Timings{}
			return err
		}

		if t.Status == pgtype.Present {
			*out = parseTimings(t.String)
		} else {
			*out = &cases.Timings{}
		}

		return nil
	})
}

var timingRegexp = regexp.MustCompile(`(?i)(\d+d)?(\d+h)?(\d+m)?`)

func parseTimings(s string) *cases.Timings {
	t := &cases.Timings{}
	if s == "" {
		return t
	}

	matches := timingRegexp.FindStringSubmatch(s)
	for _, match := range matches[1:] {
		if match == "" {
			continue
		}
		if strings.HasSuffix(match, "d") {
			t.Dd, _ = strconv.ParseInt(strings.TrimSuffix(match, "d"), 10, 64)
		} else if strings.HasSuffix(match, "h") {
			t.Hh, _ = strconv.ParseInt(strings.TrimSuffix(match, "h"), 10, 64)
		} else if strings.HasSuffix(match, "m") {
			t.Mm, _ = strconv.ParseInt(strings.TrimSuffix(match, "m"), 10, 64)
		}
	}
	return t
}
