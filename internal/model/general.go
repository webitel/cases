package model

type GeneralLookup struct {
	Id   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

func (g *GeneralLookup) SetId(id int) {
	if g.Id == nil {
		g.Id = &id
	} else {
		*g.Id = id
	}
}

func (g *GeneralLookup) GetId() *int {
	if g == nil {
		return nil
	}
	if g.Id == nil {
		return nil
	}
	return g.Id
}

func (g *GeneralLookup) SetName(name string) {
	if g.Name == nil {
		g.Name = &name
	} else {
		*g.Name = name
	}
}

func (g *GeneralLookup) GetName() *string {
	if g == nil {
		return nil
	}
	if g.Name == nil {
		return nil
	}
	return g.Name
}

type GeneralExtendedLookup struct {
	Id   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

func (g *GeneralExtendedLookup) SetId(id int) {
	if g.Id == nil {
		g.Id = &id
	} else {
		*g.Id = id
	}
}

func (g *GeneralExtendedLookup) GetId() *int {
	if g == nil {
		return nil
	}
	if g.Id == nil {
		return nil
	}
	return g.Id
}

func (g *GeneralExtendedLookup) SetName(name string) {
	if g.Name == nil {
		g.Name = &name
	} else {
		*g.Name = name
	}
}

func (g *GeneralExtendedLookup) GetName() *string {
	if g == nil {
		return nil
	}
	if g.Name == nil {
		return nil
	}
	return g.Name
}

func (g *GeneralExtendedLookup) SetType(t string) {
	if g.Type == nil {
		g.Type = &t
	} else {
		*g.Type = t
	}
}

func (g *GeneralExtendedLookup) GetType() *string {
	if g.Type == nil {
		return nil
	}
	return g.Type
}
