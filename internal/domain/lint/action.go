package lint

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

func (a *Action) Lint() string {
	// do all the fun stuff here
	services, err := a.repository.All()
	if err != nil {
		log.Println(err)
	}

	var output []string

	for _, s := range services {
		output = append(output, fmt.Sprintf("* %s", s.Id))

		for _, l := range s.Lints {
			output = append(output, fmt.Sprintf("*  %s", l))
		}
	}
	return strings.Join(output, "\n")
}
