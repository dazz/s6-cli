package filesystem

import (
	"fmt"
	"github.com/dazz/s6-cli/internal/domain/remove"
	"github.com/dazz/s6-cli/internal/domain/service"
	"testing"
)

func Test_Remove(t *testing.T) {
	rootPath := "../../../../examples/s6-overlay/s6-rc.d"
	repo := NewFilesystem(rootPath)

	t.Run("Oneshot", func(t *testing.T) {
		// create the oneshot we want to delete
		s := service.NewService("test-remove")
		s.Type = service.TypeOneshot

		if err := createOneshot(repo, s); err != nil {
			t.Errorf(err.Error())
		}

		// remove the oneshot
		steps := remove.OneshotSteps(repo, s)
		if err := remove.Remove(steps); err != nil {
			t.Errorf(err.Error())
		}

		servicePath := rootPath + "/test-remove"
		serviceScriptFilePath, _ := repo.ServiceScriptFilePath(s.Id)

		if repo.AnyFileExists([]string{
			servicePath,
			servicePath + "/dependencies.d",
			servicePath + "/dependencies.d/base",
			servicePath + "/type",
			servicePath + "/up",
			serviceScriptFilePath,
		}) {
			t.Errorf("Service was not deleted")
		}
	})

	t.Run("longrun", func(t *testing.T) {
		s := service.NewService("test-remove")
		s.Type = service.TypeLongrun

		if err := createLongrun(repo, s); err != nil {
			t.Errorf(err.Error())
		}

		// remove the service
		steps := remove.LongrunSteps(repo, s)
		if err := remove.Remove(steps); err != nil {
			t.Errorf(err.Error())
		}

		servicePath := rootPath + "/test-remove"

		if repo.AnyFileExists([]string{
			servicePath,
			servicePath + "/dependencies.d",
			servicePath + "/dependencies.d/base",
			servicePath + "/type",
			servicePath + "/run",
		}) {
			t.Errorf("Service was not deleted")
		}
	})

	t.Run("bundle", func(t *testing.T) {
		s := service.NewService("test-remove")
		s.Type = service.TypeLongrun

		if err := createBundle(repo, s); err != nil {
			t.Errorf(err.Error())
		}

		// remove the service
		steps := remove.BundleSteps(repo, s)
		if err := remove.Remove(steps); err != nil {
			t.Errorf(err.Error())
		}

		servicePath := rootPath + "/test-remove"

		if repo.AnyFileExists([]string{
			servicePath,
			servicePath + "/contents.d",
			servicePath + "/type",
		}) {
			t.Errorf("Service was not deleted")
		}
	})
}

func createOneshot(fs *Filesystem, s *service.Service) error {
	serviceDependenciesPath, _ := fs.ServiceDependenciesPath(s)
	serviceScriptFilePath, _ := fs.ServiceScriptFilePath(s.Id)
	_ = fs.CreateDirectory(fs.ServicePath(s.Id), 0755)
	_ = fs.CreateDirectory(serviceDependenciesPath, 0755)
	_ = fs.CreateFile(serviceScriptFilePath, []byte("#!/bin/sh\necho \"hello\"\n"), 0755)
	_ = fs.CreateTypeFile(s)
	_ = fs.CreateFile(fs.ServicePath(s.Id)+"/up", []byte(serviceScriptFilePath+"\n"), 0644)
	return nil
}

func createLongrun(fs *Filesystem, s *service.Service) error {
	serviceDependenciesPath, _ := fs.ServiceDependenciesPath(s)
	_ = fs.CreateDirectory(fs.ServicePath(s.Id), 0755)
	_ = fs.CreateDirectory(serviceDependenciesPath, 0755)
	_ = fs.CreateTypeFile(s)
	_ = fs.CreateFile(fs.ServicePath(s.Id)+"/run", []byte(fmt.Sprintf("#!/bin/sh\n\nexec %s\n", "sleep 10000")), 0644)
	return nil
}

func createBundle(fs *Filesystem, s *service.Service) error {
	serviceDependenciesPath, _ := fs.ServiceDependenciesPath(s)
	_ = fs.CreateDirectory(fs.ServicePath(s.Id), 0755)
	_ = fs.CreateDirectory(serviceDependenciesPath, 0755)
	_ = fs.CreateTypeFile(s)
	return nil
}
