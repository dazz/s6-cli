package mock

import (
	"github.com/dazz/s6-cli/internal/domain/service"
)

type Repository struct {
	byId   map[service.Id]*service.Service
	allIds []service.Id
}

func NewRepository(services []*service.Service) *Repository {
	r := &Repository{
		byId:   make(map[service.Id]*service.Service),
		allIds: []service.Id{},
	}
	for _, s := range services {
		r.byId[s.Id] = s
		r.allIds = append(r.allIds, s.Id)
	}
	return r
}

func (r Repository) CreateDirectory(path string, perm int) error {
	return nil
}

func (r Repository) All() ([]*service.Service, error) {
	var services []*service.Service
	for _, id := range r.allIds {
		services = append(services, r.byId[id])
	}
	return services, nil
}

func (r Repository) One(id service.Id) (*service.Service, error) {
	_, ok := r.byId[id]
	if !ok {
		return nil, nil
	}
	return r.byId[id], nil
}

func (r Repository) ServiceScriptFilePath(id service.Id) (string, error) {
	return "", nil
}

func (r Repository) ServicePath(id service.Id) string {
	return ""
}

func (r Repository) ServiceDependenciesPath(s *service.Service) (string, error) {
	return "", nil
}

func (r Repository) FileExists(file string) bool {
	return false
}

func (r Repository) AllFileExists(files []string) bool {
	return true
}

func (r Repository) AnyFileExists(files []string) bool {
	return false
}

func (r Repository) CreateDependencies(service *service.Service) error {
	return nil
}

func (r Repository) CreateTypeFile(service *service.Service) error {
	return nil
}

func (r Repository) CreateFile(path string, content []byte, perm int) error {
	return nil
}

func (r Repository) RemoveDirectory(path string) error {
	return nil
}

func (r Repository) RemoveFile(path string) error {
	return nil
}

func (r Repository) RemoveDependencies(service *service.Service) error {
	return nil
}

func (r Repository) FindDependenciesById(id service.Id) ([]string, error) {
	return []string{}, nil
}
