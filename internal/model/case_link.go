package model

import "time"

type CaseLink struct {
	*Author
	*Editor
	*Contact
    Id        int64      `json:"id" db:"id"`
    Ver       int32      `json:"ver" db:"ver"`
    Etag      string     `json:"etag" db:"etag"`
    CreatedAt *time.Time `json:"created_at" db:"created_at"`
    UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
    Name      *string     `json:"name" db:"name"`
    Url       string     `json:"url" db:"url"`
}