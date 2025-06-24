package model

type Calendar struct {
	Id   *int  `db:"calendar_id"`
	Name *string `db:"calendar_name"`
}

func (c *Calendar) GetId() *int {
    if c == nil {
        return nil
    }
    return c.Id
}

func (c *Calendar) GetName() *string {
    if c == nil {
        return nil
    }
    return c.Name
} 

func (c *Calendar) SetId(id int) {
	if c == nil {
		return
	}
	c.Id = &id
}

func (c *Calendar) SetName(name string) {
	if c == nil {
		return
	}
	c.Name = &name
}
