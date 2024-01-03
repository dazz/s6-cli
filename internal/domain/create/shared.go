package create

import (
	"errors"
	"github.com/dazz/s6-cli/internal/domain/service"
)

func checkServiceDoesNotExist(repo service.Repository, service *service.Service) error {
	if repo.FileExists(repo.ServicePath(service.Id)) {
		return errors.New("service already exists, remove service before creating new")
	}
	return nil
}

func createServicePath(repo service.Repository, service *service.Service) error {
	return repo.CreateDirectory(repo.ServicePath(service.Id), 0755)
}

func createTypeFile(repo service.Repository, service *service.Service) error {
	return repo.CreateTypeFile(service)
}

func createDependencies(repo service.Repository, service *service.Service) error {
	dependenciesPath, err := repo.ServiceDependenciesPath(service)
	if err != nil {
		return err
	}
	if err := repo.CreateDirectory(dependenciesPath, 0755); err != nil {
		return err
	}
	if err := repo.CreateDependencies(service); err != nil {
		return err
	}
	return nil
}
