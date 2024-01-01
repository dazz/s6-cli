package create

import (
	"github.com/dazz/s6-cli/internal/domain/service"
	"github.com/dazz/s6-cli/internal/infrastructure/persistence/mock"
	"testing"
)

func Test_Command_Execute(t *testing.T) {

	t.Run("returns err for none existing service", func(t *testing.T) {
		repo := &mock.Repository{}
		c := NewCommand(repo, service.Id("test"), "not-existing-type")
		_, err := c.Execute()

		if err != nil {
			t.Errorf("Execute() must return an error")
		}
	})

	t.Run("returns id of created service", func(t *testing.T) {
		repo := &mock.Repository{}
		c := NewCommand(repo, "test", "o")
		result, _ := c.Execute()

		if result != "test" {
			t.Errorf("Execute() must return the created id")
		}
	})
}
