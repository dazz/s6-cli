package create

import (
	"fmt"
	"github.com/dazz/s6-cli/internal/domain/service"
)

func LongrunSteps(repo service.Repository, s *service.Service) service.StepCollection {
	return service.StepCollection{
		Repo:    repo,
		Service: s,
		Steps: []func(repo service.Repository, s *service.Service) error{
			checkServiceDoesNotExist,
			createServicePath,
			createTypeFile,
			createLongrunServiceFile,
			createDependencies,
		},
	}
}

func createLongrunServiceFile(repo service.Repository, service *service.Service) error {
	// Iterate the up file
	runFile := repo.ServicePath(service.Id) + "/run"
	runFileContent := []byte(fmt.Sprintf("#!/bin/sh\n\nexec %s\n", "sleep 10000"))

	return repo.CreateFile(runFile, runFileContent, 0644)
}
