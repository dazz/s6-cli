package remove

import (
	"errors"
	"github.com/dazz/s6-cli/internal/domain/service"
)

type Command struct {
	repository service.Repository
	id         service.Id
}

func NewCommand(repository service.Repository, id service.Id) *Command {
	return &Command{
		repository: repository,
		id:         id,
	}
}

// Execute does all the fun stuff
func (c *Command) Execute() (string, error) {

	s, err := c.repository.One(c.id)
	if err != nil {
		return "", err
	}

	var steps service.StepIterator

	switch s.Type {
	case service.TypeOneshot:
		steps = OneshotSteps(c.repository, s)
	case service.TypeLongrun:
		steps = LongrunSteps(c.repository, s)
	case service.TypeBundle:
		steps = BundleSteps(c.repository, s)
	default:
		return "", errors.New("unknown service type")
	}

	if err := Remove(steps); err != nil {
		return "", err
	}

	return string(c.id), nil
}
