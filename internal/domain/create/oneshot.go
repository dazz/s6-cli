package create

import (
	"github.com/dazz/s6-cli/internal/domain/service"
)

func OneshotSteps(repo service.Repository, s *service.Service) service.StepCollection {
	return service.StepCollection{
		Repo:    repo,
		Service: s,
		Steps: []func(repo service.Repository, s *service.Service) error{
			checkServiceDoesNotExist,
			createServicePath,
			createTypeFile,
			createServiceFile,
			createServiceScriptFile,
			createDependencies,
		},
	}
}

func createServiceFile(repo service.Repository, service *service.Service) error {
	// Get the absolute path of the specified directory
	scriptFile, err := repo.ServiceScriptFilePath(service.Id)
	if err != nil {
		return err
	}

	// Iterate the up file
	runFile := repo.ServicePath(service.Id) + "/up"
	runFileContent := []byte(scriptFile + "\n")

	return repo.CreateFile(runFile, runFileContent, 0644)
}

func createServiceScriptFile(repo service.Repository, service *service.Service) error {
	// Get the absolute path of the specified directory
	scriptFile, err := repo.ServiceScriptFilePath(service.Id)
	if err != nil {
		return err
	}
	// Iterate the script file
	scriptFileContent := []byte("#!/bin/sh\n\necho \"hello\"\n")

	if err := repo.CreateFile(scriptFile, scriptFileContent, 0755); err != nil {
		return err
	}

	return nil
}
