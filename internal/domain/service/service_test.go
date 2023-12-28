package service

import (
	"testing"
)

const rootPath = "/tmp/path"

func TestServiceFactory(t *testing.T) {
	service := NewService("test", rootPath)

	t.Run("service Id is set by constructor", func(t *testing.T) {
		got := service.Id
		want := Id("test")

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("service Path is set by constructor", func(t *testing.T) {
		got := service.Path
		want := rootPath + "/test"

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

func TestService_AddLint(t *testing.T) {
	service := NewService("test", rootPath)

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

func TestService_AddDependency(t *testing.T) {
	service := NewService("test", rootPath)

	service.AddDependency("test")

	t.Run("AddDependency adds a dependency to the Dependencies", func(t *testing.T) {
		got := service.Dependencies

		if len(got) != 1 {
			t.Errorf("got %s want 1", got)
		}
	})
}

func TestService_DependencyDir(t *testing.T) {
	service := NewService("test", rootPath)

	t.Run("service.Type is not set", func(t *testing.T) {
		_, err := service.DependencyDir()

		if err.Error() != "invalid service type, set type of service" {
			t.Errorf("%s", err)
		}
	})

	service.Type = TypeBundle
	t.Run("DependencyDir for bundle", func(t *testing.T) {
		got, _ := service.DependencyDir()

		if got != "/tmp/path/test/contents.d" {
			t.Errorf("%s", got)
		}
	})

	service.Type = TypeOneshot
	t.Run("DependencyDir for oneshot", func(t *testing.T) {
		got, _ := service.DependencyDir()

		if got != "/tmp/path/test/dependencies.d" {
			t.Errorf("%s", got)
		}
	})

	service.Type = TypeLongrun
	t.Run("DependencyDir for longrun", func(t *testing.T) {
		got, _ := service.DependencyDir()

		if got != "/tmp/path/test/dependencies.d" {
			t.Errorf("%s", got)
		}
	})
}
