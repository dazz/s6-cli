package service

import (
	"testing"
)

const rootPath = "/tmp/path"

func Test_ServiceFactory(t *testing.T) {
	service := NewService("test")

	t.Run("service Id is set by constructor", func(t *testing.T) {
		got := service.Id
		want := Id("test")

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("service Lints is set by constructor", func(t *testing.T) {
		got := service.Lints

		if len(got) != 0 {
			t.Errorf("got %s want 0", got)
		}
	})
}

func Test_Service_AddLint(t *testing.T) {
	service := NewService("test")

	service.AddLint("test")

	t.Run("AddLint adds a lint to the Lints", func(t *testing.T) {
		got := service.Lints

		if len(got) != 1 {
			t.Errorf("got %s want 1", got)
		}
	})

	t.Run("AddLint does not add empty lint to the Lints", func(t *testing.T) {
		service.AddLint("")

		got := service.Lints

		if len(got) != 1 {
			t.Errorf("got %s want 1", got)
		}
	})
}

func Test_Service_AddDependency(t *testing.T) {
	s := NewService("test")

	t.Run("adds a dependency to the createDependencies", func(t *testing.T) {
		s.AddDependency("test")

		if got := s.Dependencies; len(got) != 1 {
			t.Errorf("got %s want 1", got)
		}
	})
}
