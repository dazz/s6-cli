package filesystem

import (
	"github.com/dazz/s6-cli/internal/domain/service"
	"path/filepath"
	"testing"
)

func Test_ServiceScriptFilePath(t *testing.T) {
	rootPath := "../../../../examples/s6-overlay/s6-rc.d"
	repo := NewFilesystem(rootPath)
	s := service.NewService("test")

	t.Run("returns absolute path to script file", func(t *testing.T) {
		want, _ := filepath.Abs(rootPath + "/../scripts")
		if got, _ := repo.ServiceScriptFilePath(s.Id); got != want+"/test" {
			t.Errorf("%s but wanted %s", got, want)
		}
	})
}

func Test_ServiceDependencyPath(t *testing.T) {
	rootPath := "../../../../examples/s6-overlay/s6-rc.d"
	repo := NewFilesystem(rootPath)
	s := service.NewService("test")

	t.Run("Type is not set", func(t *testing.T) {
		_, err := repo.ServiceDependencyPath(s)

		if err.Error() != "invalid service type, set type of service" {
			t.Errorf("%s", err)
		}
	})

	t.Run("for bundle", func(t *testing.T) {
		s.Type = service.TypeBundle
		got, _ := repo.ServiceDependencyPath(s)

		if got != "../../../../examples/s6-overlay/s6-rc.d/test/contents.d" {
			t.Errorf("%s", got)
		}
	})

	t.Run("for oneshot", func(t *testing.T) {
		s.Type = service.TypeOneshot
		got, _ := repo.ServiceDependencyPath(s)

		if got != "../../../../examples/s6-overlay/s6-rc.d/test/dependencies.d" {
			t.Errorf("%s", got)
		}
	})

	t.Run("for longrun", func(t *testing.T) {
		s.Type = service.TypeLongrun
		got, _ := repo.ServiceDependencyPath(s)

		if got != "../../../../examples/s6-overlay/s6-rc.d/test/dependencies.d" {
			t.Errorf("%s", got)
		}
	})
}

func Test_ServicePath(t *testing.T) {
	s := service.NewService("test")
	rootPath := "../../../../examples/s6-overlay/s6-rc.d"
	repo := NewFilesystem(rootPath)

	t.Run("returns rootPath when '' is passed", func(t *testing.T) {
		if got := repo.ServicePath(""); got != rootPath {
			t.Errorf("%s", got)
		}
	})

	t.Run("returns path when s.id is passed", func(t *testing.T) {
		if got := repo.ServicePath(s.Id); got != "../../../../examples/s6-overlay/s6-rc.d/test" {
			t.Errorf("%s", got)
		}
	})
}
