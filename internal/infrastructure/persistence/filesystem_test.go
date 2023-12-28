package persistence

import (
	"testing"
)

const rootPath = "../../../examples/s6-overlay/s6-rc.d"

func TestFilesystemFactory(t *testing.T) {
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
