package remove

import (
	"github.com/dazz/s6-cli/internal/domain/service"
	"github.com/dazz/s6-cli/internal/infrastructure/persistence/mock"
	"testing"
)

func Test_Command_Execute(t *testing.T) {

	t.Run("Execute() returns error", func(t *testing.T) {

		repo := mock.NewRepository([]*service.Service{
			service.NewService("test"),
		})
		c := NewCommand(repo, "test")
		_, err := c.Execute()
		if err == nil {
			t.Errorf("Execute() must return an error.")
		}
		if err.Error() != "unknown service type" {
			t.Errorf("Execute() must return an error with message 'unknown service type'. Error: %s", err.Error())
		}
	})

	t.Run("Execute() returns string", func(t *testing.T) {
		s := service.NewService("test")
		s.Type = service.TypeOneshot

		repo := mock.NewRepository([]*service.Service{
			s,
		})
		c := NewCommand(repo, "test")
		result, _ := c.Execute()

		if result == "" {
			t.Errorf("Execute() must return a string")
		}
	})
}
