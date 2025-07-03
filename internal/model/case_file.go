package model

import "time"

type CaseFile struct {
	*Author
	*Contact
    Id        int       `db:"id"`
    CreatedAt *time.Time `db:"created_at"`
    Size      int64     `db:"size"`
    Mime      string    `db:"mime"`
    Name      string    `db:"name"`
    Url       string    `db:"url"`
}