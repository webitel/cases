package model

type Contact struct {
    Id   *int    `db:"contact_id"`
    Name *string `db:"contact_name"`
}

func (a *Contact) GetId() *int {
	if a == nil {
		return nil
	}
	return a.Id
}

func (a *Contact) GetName() *string {
	if a == nil {
		return nil
	}
	return a.Name
}

func (a *Contact) SetId(id int) {
	if a == nil {
		return
	}
	a.Id = &id
}

func (a *Contact) SetName(name string) {
	if a == nil {
		return
	}
	a.Name = &name
}
