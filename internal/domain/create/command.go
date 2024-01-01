package create

import (
	"github.com/dazz/s6-cli/internal/domain/service"
)

type Command struct {
	repository  service.Repository
	id          service.Id
	serviceType service.Type
}

func NewCommand(repository service.Repository, id service.Id, serviceType service.Type) *Command {
	return &Command{
		repository:  repository,
		serviceType: serviceType,
		id:          id,
	}
}

func (c *Command) Execute() (string, error) {

	s := service.NewService(c.id)
	switch c.serviceType {
	case service.TypeOneshot:
		s.Type = service.TypeOneshot
		if err := c.repository.Oneshot(s); err != nil {
			return "", err
		}
	case service.TypeLongrun:
		s.Type = service.TypeLongrun
		if err := c.repository.Longrun(s); err != nil {
			return "", err
		}
	case service.TypeBundle:
		s.Type = service.TypeLongrun
		if err := c.repository.Bundle(s); err != nil {
			return "", err
		}
	}

	return string(c.id), nil
}
