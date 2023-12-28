package lint

import (
	"github.com/dazz/s6-cli/internal/infrastructure/persistence/mock"
	"testing"
)

func TestAction_Lint(t *testing.T) {
	t.Run("Lint() returns string", func(t *testing.T) {
		repo := &mock.Repository{}
		action := NewAction(repo)
		result := action.Lint()

		if result == "" {
			t.Errorf("Lint() must return a string")
		}
	})
}
