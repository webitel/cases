package model

type Author struct {
	Id   *int64  `db:"created_by_id"`   // or updated_by_id
	Name *string `db:"created_by_name"` // or updated_by_name
}

type Editor struct {
	Id   *int64  `db:"updated_by_id"`   // or created_by_id
	Name *string `db:"updated_by_name"` // or created_by_name
}
