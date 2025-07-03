package model

type Editor struct {
	Id   *int    `db:"updated_by_id"`
	Name *string `db:"updated_by_name"`
}

func (a *Editor) GetId() *int {
	if a == nil {
		return nil
	}
	return a.Id
}

func (a *Editor) GetName() *string {
	if a == nil {
		return nil
	}
	return a.Name
}

func (a *Editor) SetId(id int) {
	if a == nil {
		return
	}
	a.Id = &id
}

func (a *Editor) SetName(name string) {
	if a == nil {
		return
	}
	a.Name = &name
}
