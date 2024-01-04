package filesystem

import (
	"errors"
	"github.com/dazz/s6-cli/internal/domain/service"
	"os"
)

func (fs *Filesystem) CreateFile(path string, content []byte, perm int) error {
	return os.WriteFile(path, content, os.FileMode(perm))
}

func (fs *Filesystem) CreateTypeFile(service *service.Service) error {
	// Iterate the type file
	typeFile := fs.ServicePath(service.Id) + "/type"
	typeFileContent := []byte(string(service.Type) + "\n")

	return fs.CreateFile(typeFile, typeFileContent, 0644)
}

func (fs *Filesystem) CreateDirectory(path string, perm int) error {

	if err := os.Mkdir(path, os.FileMode(perm)); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func (fs *Filesystem) CreateDependencies(service *service.Service) error {
	dependenciesPath, err := fs.ServiceDependenciesPath(service)
	if err != nil {
		return err
	}
	if !fs.FileExists(dependenciesPath) {
		return errors.New("dependencies path does not exist")
	}

	// Iterate the dependencies
	for _, dependency := range service.Dependencies {
		// Iterate the type file
		depFile := dependenciesPath + "/" + string(dependency)
		depFileContent := []byte("")

		if err := fs.CreateFile(depFile, depFileContent, 0644); err != nil {
			return err
		}
	}
	return nil
}
