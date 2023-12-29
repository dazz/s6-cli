package mermaid

import (
	"github.com/dazz/s6-cli/internal/infrastructure/persistence/mock"
	"testing"
)

func TestAction_Lint(t *testing.T) {
	t.Run("Output() returns string", func(t *testing.T) {
		repo := &mock.Repository{}
		action := NewAction(repo)
		result := action.Output()

		if result == "" {
			t.Errorf("Output() must return a string")
		}
	})
}
