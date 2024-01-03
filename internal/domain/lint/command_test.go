package lint

import (
	"github.com/dazz/s6-cli/internal/infrastructure/persistence/mock"
	"testing"
)

func Test_Command_Execute(t *testing.T) {
	t.Run("Execute() returns string", func(t *testing.T) {
		repo := &mock.Repository{}
		c := NewCommand(repo)
		result, err := c.Execute()

		if err != nil {
			t.Errorf("Execute() must not return an error. Error: %s", err.Error())
		}

		if result != "" {
			t.Errorf("Execute() must return a string, returned %q", result)
		}
	})
}
