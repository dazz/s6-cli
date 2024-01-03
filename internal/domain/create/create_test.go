package create

import (
	"github.com/dazz/s6-cli/internal/domain/service"
	"github.com/dazz/s6-cli/internal/infrastructure/persistence/mock"
	"testing"
)

func Test_Create(t *testing.T) {
	repo := mock.Repository{}
	s := service.NewService("test")
	t.Run("Oneshot service", func(t *testing.T) {
		// create the oneshot we want to delete
		withSteps := service.StepCollection{
			Repo:    repo,
			Service: s,
			Steps: []func(repo service.Repository, s *service.Service) error{
				func(repo service.Repository, s *service.Service) error {
					return nil
				},
			},
		}
		err := Create(withSteps)
		if err != nil {
			t.Errorf("StepCollection returned an error: %s", err)
		}
	})
}
