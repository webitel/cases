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

// TimeStringToTimestamp converts a time string into a Unix timestamp in milliseconds.
func TimeStringToTimestamp(timeStr string, format ...string) (int64, error) {
	// If no custom format is provided, use the default format.
	timeFormat := TimeStamp
	if len(format) > 0 {
		timeFormat = format[0]
	}

	// Check if the time string has a timezone. If not, append 'Z' for UTC.
	if !hasTimeZone(timeStr) {
		timeStr += "Z"
	}

	// Parse the time string using the provided or default format.
	parsedTime, err := time.Parse(timeFormat, timeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse time string: %v", err)
	}

	// Convert the parsed time to a Unix timestamp in milliseconds.
	return Timestamp(parsedTime), nil
}

// hasTimeZone checks if a time string contains a timezone indicator.
func hasTimeZone(timeStr string) bool {
	return len(timeStr) > 10 && (timeStr[len(timeStr)-1] == 'Z' || timeStr[len(timeStr)-6] == '+' || timeStr[len(timeStr)-6] == '-')
}

// LocalTime returns the local Time corresponding to the given Unix time,
// possibly precised with UnixToTimestamp. Returns nil if ts is 0.
func LocalTime(ts int64) *time.Time {
	if ts > 0 {
		t := time.UnixMilli(ts).UTC()
		return &t
	}
	return nil // Return nil if `ts` is 0 or negative
}
