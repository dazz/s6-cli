package mermaid

import (
	"github.com/dazz/s6-cli/internal/infrastructure/persistence/mock"
	"testing"
)

func Test_Command_Execute(t *testing.T) {
	t.Run("Execute() returns string", func(t *testing.T) {
		repo := &mock.Repository{}
		command := NewCommand(repo)
		result, _ := command.Execute()

		if result == "" {
			t.Errorf("Execute() must return a string")
		}
	})
}
