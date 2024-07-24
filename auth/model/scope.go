package model

type Scope struct {
	Id     int64
	Name   string
	Abac   bool
	Obac   bool
	Rbac   bool
	Class  string
	Access string
}

func (s *Scope) IsRbacUsed() bool {
	return s.IsObacUsed() && s.Rbac
}

func (s *Scope) IsObacUsed() bool {
	return s.Obac
}
