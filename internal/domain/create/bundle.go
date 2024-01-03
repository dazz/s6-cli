package create

import "github.com/dazz/s6-cli/internal/domain/service"

func BundleSteps(repo service.Repository, s *service.Service) service.StepCollection {
	return service.StepCollection{
		Repo:    repo,
		Service: s,
		Steps: []func(repo service.Repository, s *service.Service) error{
			checkServiceDoesNotExist,
			createServicePath,
			createTypeFile,
			createDependencies,
		},
	}
}
