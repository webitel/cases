package model

import "time"

type DeleteOptions struct {
	IDs  []int64
	Time time.Time
}
