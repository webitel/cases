package util

import (
	"fmt"
	"time"
)

// TimeStamp default string format
const TimeStamp = `2006-01-02 15:04:05.000`

var unixepoch = time.Unix(0, 0).UTC()

// TimestampPrecision defines defaults timestamp precition (time.Millisecond)
// const TimestampPrecision int64 = 1e6 // time.Millisecond
// const (
// 	unixToInternal = 1e3 // sec*ms
// 	internalToUnix = 1e6 // nano/ms
// )

// UnixToTimestamp defines default timestamp precition
// TimestampToUnix = (time.Second / UnixToTimestamp) // 1e3
const UnixToTimestamp = time.Millisecond // milliseconds // 1e6

// func SetTimestampPrecision( time.Duration)

// Timestamp returns number of seconds,
// posibly precised with UnixToTimestamp,
// elapsed since January 1, 1970 UTC
// until given at local time
func Timestamp(at time.Time) (ts int64) {
	if at.IsZero() || at.Before(unixepoch) {
		return 0
	}
	switch UnixToTimestamp {
	case time.Second:
		return at.Unix() // seconds
	case time.Nanosecond:
		return at.UnixNano()
	case time.Millisecond, time.Microsecond:
		return at.UnixNano() / (int64)(UnixToTimestamp)
	default:
		panic(fmt.Errorf(`timestamp: invalid precision %s`, UnixToTimestamp))
	}
}

// LocalTime returns the local Time
// corresponding to the given Unix time,
// posibly precised with UnixToTimestamp
func LocalTime(ts int64) (at time.Time) {
	if ts > 0 {
		timestampToUnix := (int64)(time.Second / UnixToTimestamp)                       // time.Second(1e9) / time.Millicesond(1e6) = 1e3
		at = time.Unix(ts/timestampToUnix, ts%timestampToUnix*(int64)(UnixToTimestamp)) // *1e9) // *time.Second
	}
	return // at.IsZero()
}
