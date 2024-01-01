package filesystem

import (
	"github.com/dazz/s6-cli/internal/domain/service"
	"testing"
)

func Test_Remove(t *testing.T) {
	rootPath := "../../../../examples/s6-overlay/s6-rc.d"
	repo := NewFilesystem(rootPath)
	s := service.NewService("test")

	t.Run("removes the oneshot service including the script file", func(t *testing.T) {
		// create the oneshot we want to delete
		_ = repo.Oneshot(s)
		s.Type = service.TypeOneshot
		if err := repo.Remove(s); err != nil {
			t.Errorf(err.Error())
		}
		if err := repo.fileExists(rootPath + "/test"); err == nil {
			t.Errorf("Service was not deleted")
		}
	})

	t.Run("removes the oneshot service including the script file", func(t *testing.T) {
		// create the oneshot we want to delete
		_ = repo.Longrun(s)
		s.Type = service.TypeLongrun

		if err := repo.Remove(s); err != nil {
			t.Errorf(err.Error())
		}

		if err := repo.fileExists(rootPath + "/test"); err == nil {
			t.Errorf("Service was not deleted")
		}
	})
}
