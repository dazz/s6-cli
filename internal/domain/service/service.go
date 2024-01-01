package service

const (
	TypeOneshot Type = "oneshot"
	TypeLongrun Type = "longrun"
	TypeBundle  Type = "bundle"
)

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
