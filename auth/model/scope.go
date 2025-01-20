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
	if s == nil {
		return false
	}
	return s.Rbac
}

func (s *Scope) IsObacUsed() bool {
	if s == nil {
		return false
	}
	return s.Obac
}
