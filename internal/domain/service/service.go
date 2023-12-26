package service

import (
	"errors"
)

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

func NewService(id Id, rootPath string) *Service {
	return &Service{
		Id:    id,
		Path:  rootPath + "/" + string(id),
		Valid: true,
		Lints: []Lint{},
	}
}

func (s *Service) AddLint(lint string) {
	s.Lints = append(s.Lints, Lint(lint))
	s.Valid = false
}

func (s *Service) DependencyDir() (string, error) {
	if s.Type == TypeBundle {
		return s.Path + "/contents.d", nil
	}
	if s.Type == TypeOneshot || s.Type == TypeLongrun {
		return s.Path + "/dependencies.d", nil
	}
	return "", errors.New("invalid service type")
}

func (s *Service) AddDependency(dependency string) {
	s.Dependencies = append(s.Dependencies, Id(dependency))
}
