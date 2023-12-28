package mock

import "github.com/dazz/s6-cli/internal/domain/service"

type Repository struct {
}

func (r Repository) All() ([]*service.Service, error) {
	return []*service.Service{
		service.NewService("abc", "/tmp/path"),
	}, nil
}
