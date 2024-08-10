package model

type Scope struct {
	Name   string
	Class  string
	Access string
	Id     int64
	Abac   bool
	Obac   bool
	Rbac   bool
}

func (s *Scope) IsRbacUsed() bool {
	return s.IsObacUsed() && s.Rbac
}

func (s *Scope) IsObacUsed() bool {
	return s.Obac
}
