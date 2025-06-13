package model

type Author struct {
	Id   *int    `db:"created_by_id"`
	Name *string `db:"created_by_name"`
}

func (a *Author) SetId(id int) {
	if a == nil {
		return
	}
	a.Id = &id
}

func (a *Author) SetName(name string) {
	if a == nil {
		return
	}
	a.Name = &name
}

func (a *Author) GetId() *int {
	if a == nil {
		return nil
	}
	return a.Id
}

func (a *Author) GetName() *string {
	if a == nil {
		return nil
	}
	return a.Name
}
