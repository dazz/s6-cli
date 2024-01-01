package remove

import (
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

	if err := c.repository.Remove(s); err != nil {
		return "", err
	}
	return string(c.id), nil
}
