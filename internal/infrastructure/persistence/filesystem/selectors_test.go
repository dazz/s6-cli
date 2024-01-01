package filesystem

import (
	"github.com/dazz/s6-cli/internal/domain/service"
	"testing"
)

const rootPath = "../../../../examples/s6-overlay/s6-rc.d"

func Test_FilesystemFactory(t *testing.T) {
	fs := NewFilesystem(rootPath)

	t.Run("Factory gives valid instance of Filesystem", func(t *testing.T) {
		if fs == nil {
			t.Errorf("Filesystem instance must not be empty")
		}
	})
}

func Test_All(t *testing.T) {
	fs := NewFilesystem(rootPath)
	t.Run("All() throws no error", func(t *testing.T) {
		_, err := fs.All()
		if err != nil {
			t.Errorf("Error: %s", err)
		}
	})

	t.Run("All() returns services", func(t *testing.T) {
		services, _ := fs.All()
		if len(services) == 0 {
			t.Errorf("All() must return empty array")
		}
	})
}

func Test_One(t *testing.T) {
	fs := NewFilesystem(rootPath)
	t.Run("One() throws no error", func(t *testing.T) {
		_, err := fs.One(service.Id(""))
		if err != nil {
			t.Errorf("Error: %s", err)
		}
	})

	t.Run("One() without argument returns no service", func(t *testing.T) {
		s, _ := fs.One("")
		if s != nil {
			t.Errorf("One() must not be nil")
		}
	})

	t.Run("One() with argument returns service", func(t *testing.T) {
		s, _ := fs.One("user")
		if s == nil {
			t.Errorf("One() must not be nil")
		}
	})
}
