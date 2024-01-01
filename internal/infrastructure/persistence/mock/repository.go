package mock

import "github.com/dazz/s6-cli/internal/domain/service"

type Repository struct {
}

func (r Repository) All() ([]*service.Service, error) {
	return []*service.Service{
		service.NewService(service.Id("abc")),
	}, nil
}

func (r Repository) One(id service.Id) (*service.Service, error) {
	return service.NewService(service.Id("abc")), nil
}

func (r Repository) Oneshot(service *service.Service) error {
	return nil
}

func (r Repository) Longrun(service *service.Service) error {
	return nil
}

func (r Repository) Bundle(service *service.Service) error {
	return nil
}

func (r Repository) Remove(service *service.Service) error {
	return nil
}
