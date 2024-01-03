package remove

import "github.com/dazz/s6-cli/internal/domain/service"

func Remove(steps service.StepIterator) error {
	return steps.Iterate()
}

func OneshotSteps(repo service.Repository, s *service.Service) service.StepCollection {
	return service.StepCollection{
		Repo:    repo,
		Service: s,
		Steps: []func(repo service.Repository, s *service.Service) error{
			removeServicePath,
			removeDependencies,
			removeScriptFile,
		},
	}
}

func LongrunSteps(repo service.Repository, s *service.Service) service.StepCollection {
	return service.StepCollection{
		Repo:    repo,
		Service: s,
		Steps: []func(repo service.Repository, s *service.Service) error{
			removeServicePath,
			removeDependencies,
		},
	}
}

func BundleSteps(repo service.Repository, s *service.Service) service.StepCollection {
	return service.StepCollection{
		Repo:    repo,
		Service: s,
		Steps: []func(repo service.Repository, s *service.Service) error{
			removeServicePath,
			removeDependencies,
		},
	}
}

func removeServicePath(repo service.Repository, s *service.Service) error {
	return repo.RemoveDirectory(repo.ServicePath(s.Id))
}

func removeScriptFile(repo service.Repository, s *service.Service) error {
	scriptFilePath, err := repo.ServiceScriptFilePath(s.Id)
	if err != nil {
		return err
	}
	return repo.RemoveFile(scriptFilePath)
}

func removeDependencies(repo service.Repository, s *service.Service) error {
	dependencyFoundAt, err := repo.FindDependenciesById(s.Id)
	if err != nil {
		return err
	}

	if len(dependencyFoundAt) > 0 {
		for _, depFile := range dependencyFoundAt {
			if err := repo.RemoveFile(depFile); err != nil {
				return err
			}
		}
	}
	return nil
}
