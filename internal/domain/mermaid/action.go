package mermaid

import (
	"fmt"
	"github.com/dazz/s6-cli/internal/domain/service"
	"log"
	"strings"
)

type Action struct {
	repository service.Repository
}

func NewAction(repository service.Repository) *Action {
	return &Action{
		repository: repository,
	}
}

func (a *Action) Output() string {
	// do all the fun stuff hereLint
	services, err := a.repository.All()
	if err != nil {
		log.Println(err)
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
	return strings.Join(output, "\n")
}
