package service

type StepIterator interface {
	Iterate() error
}

type StepCollection struct {
	Repo    Repository
	Service *Service
	Steps   []func(repo Repository, s *Service) error
}

// Iterate implements the StepIterator interface
func (c StepCollection) Iterate() error {
	for _, step := range c.Steps {
		if err := step(c.Repo, c.Service); err != nil {
			return err
		}
	}
	return nil
}
