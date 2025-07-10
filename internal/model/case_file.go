package model

import "time"

type CaseFile struct {
	Id        int                    `db:"id"`
	CreatedAt *time.Time             `db:"created_at"`
	Size      int64                  `db:"size"`
	Mime      string                 `db:"mime"`
	Name      string                 `db:"name"`
	Url       string                 `db:"url"`
	Author    *GeneralExtendedLookup `db:"created_by"`
	Source    string                 `db:"source"`
}
