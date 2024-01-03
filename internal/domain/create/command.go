package create

import (
	"errors"
	"github.com/dazz/s6-cli/internal/domain/service"
)

const defaultDependency = "base"

type Command struct {
	repository  service.Repository
	id          service.Id
	serviceType service.Type
}

func NewCommand(repository service.Repository, id service.Id, serviceType service.Type) *Command {
	return &Command{
		repository:  repository,
		id:          id,
		serviceType: serviceType,
	}
}

func (c *Command) Execute() (string, error) {

	s := service.NewService(c.id)
	s.Type = c.serviceType
	s.Dependencies = []service.Id{defaultDependency}

	var steps service.StepIterator

	switch c.serviceType {
	case service.TypeOneshot:
		steps = OneshotSteps(c.repository, s)
	case service.TypeLongrun:
		steps = LongrunSteps(c.repository, s)
	case service.TypeBundle:
		steps = BundleSteps(c.repository, s)
	default:
		return "", errors.New("unknown service type")
	}

	if err := Create(steps); err != nil {
		return "", err
	}

	return string(c.id), nil
}
