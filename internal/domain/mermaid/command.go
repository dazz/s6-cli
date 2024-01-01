package mermaid

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

// Execute does all the fun stuff
func (a *Command) Execute() (string, error) {

	services, err := a.repository.All()
	if err != nil {
		return "", err
	}

	var output []string
	output = append(output, "```mermaid")
	output = append(output, "graph TD;")
	for _, s := range services {
		for _, dependency := range s.Dependencies {
			output = append(output, fmt.Sprintf("    %s --> %s", s.Id, dependency))
		}
	}
	output = append(output, "```")
	return strings.Join(output, "\n"), nil
}
