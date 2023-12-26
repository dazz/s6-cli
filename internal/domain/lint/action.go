package lint

import (
	"fmt"
	"github.com/dazz/s6-cli/internal/domain/service"
	"log"
)

type Action struct {
	repository service.Repository
}

func NewAction(repository service.Repository) *Action {
	return &Action{
		repository: repository,
	}
}

func (a *Action) Lint() bool {
	// do all the fun stuff here
	services, err := a.repository.All()
	if err != nil {
		log.Println(err)
	}

	for _, s := range services {
		fmt.Println("* " + s.Id)
		for _, l := range s.Lints {
			fmt.Println("  * " + l)
		}
	}
	return false
}
