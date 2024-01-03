package lint

import (
	"fmt"
	"github.com/dazz/s6-cli/internal/domain/service"
	"strings"
)

type Command struct {
	repository service.Repository
}

func NewCommand(repository service.Repository) *Command {
	return &Command{
		repository: repository,
	}
}

func (a *Command) Execute() (string, error) {
	// do all the fun stuff here
	services, err := a.repository.All()
	if err != nil {
		return "", err
	}

	var output []string

	for _, s := range services {
		if len(s.Lints) == 0 {
			continue
		}
		output = append(output, fmt.Sprintf("* %s", s.Id))
		for _, l := range s.Lints {
			output = append(output, fmt.Sprintf("  * %s", l))
		}
	}
	return strings.Join(output, "\n"), nil
}
