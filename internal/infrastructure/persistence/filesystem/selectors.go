package filesystem

import (
	"github.com/dazz/s6-cli/internal/domain/service"
	"log"
)

func (fs *Filesystem) All() ([]*service.Service, error) {
	err := fs.compile(RootService)
	if err != nil {
		log.Println(err)
	}
	var services []*service.Service
	for _, id := range fs.allIds {
		services = append(services, fs.byId[id])
	}
	return services, nil
}

func (fs *Filesystem) One(id service.Id) (*service.Service, error) {
	if id == "" {
		return nil, nil
	}

	if err := fs.compile(id); err != nil {
		return nil, err
	}
	if fs.contains(id) {
		return fs.byId[id], nil
	}
	return nil, nil
}

// All returns all services
func (fs *Filesystem) contains(id service.Id) bool {
	_, ok := fs.byId[id]
	if ok {
		return true
	}
	return false
}
