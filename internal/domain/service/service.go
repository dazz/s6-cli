// Package service contains the domain layer for the service struct.
// The service struct is the main entity of the application.
package service

const (
	TypeOneshot Type = "oneshot"
	TypeLongrun Type = "longrun"
	TypeBundle  Type = "bundle"
)

var Types = []Type{TypeOneshot, TypeLongrun, TypeBundle}

func ValidType(ty Type) bool {
	for _, t := range Types {
		if t == ty {
			return true
		}
	}
	return false
}

type Id string
type Type string
type Lint string

type Service struct {
	Id           Id
	Type         Type
	Dependencies []Id
	Lints        []Lint
	Valid        bool
	Path         string
}

func NewService(id Id) *Service {
	return &Service{
		Id:    id,
		Lints: make([]Lint, 0),
	}
}

func (s *Service) AddLint(lint string) {
	if lint == "" {
		return
	}
	s.Lints = append(s.Lints, Lint(lint))
}

func (s *Service) AddDependency(dependency string) {
	s.Dependencies = append(s.Dependencies, Id(dependency))
}
