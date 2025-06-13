package model

type Lookup interface {
	SetId(int)
	GetId() *int
	SetName(string)
	GetName() *string
}
