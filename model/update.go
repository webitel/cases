package model

import "time"

type UpdateOptions struct {
	ID     int64
	Fields []string
	Time   time.Time
}
