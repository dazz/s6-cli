package persistence

import (
	"errors"
	"fmt"
	"github.com/dazz/s6-cli/internal/domain/service"
	"log"
	"os"
	"strings"
)

const RootService service.Id = "user"

type Filesystem struct {
	rootPath string
	byId     map[service.Id]*service.Service
	allIds   []service.Id
}

func NewFilesystem(rootPath string) *Filesystem {
	return &Filesystem{
		rootPath: rootPath,
		byId:     make(map[service.Id]*service.Service),
		allIds:   []service.Id{},
	}
}

func (fs *Filesystem) All() ([]*service.Service, error) {
	err := fs.compile(RootService)
	if err != nil {
		log.Println(err)
	}
	var services []*service.Service
	for _, id := range fs.allIds {
		services = append(services, fs.byId[id])
		//fmt.Println(id)
	}
	return services, nil
}

func (fs *Filesystem) contains(id service.Id) bool {
	_, ok := fs.byId[id]
	if ok {
		return true
	}
	return false
}

// compile all folders and directories to a list of services we can work with
func (fs *Filesystem) compile(id service.Id) error {
	if id == "" {
		id = RootService
	}

	// Check if the service is already in the services list
	if fs.contains(id) {
		return nil
	}

	s := service.NewService(id, fs.rootPath)

	// add the service to the services list
	fs.byId[id] = s
	fs.allIds = append(fs.allIds, id)

	// Check if the directory exists
	if _, err := os.Stat(s.Path); os.IsNotExist(err) {
		s.AddLint(fmt.Sprintf("invalid directory: path %s does not exist", s.Path))
		return nil
	}

	serviceType, err := fs.serviceType(s)
	if err != nil {
		s.AddLint(err.Error())
		return nil
	}
	s.Type = serviceType

	// check if the run file exists
	if serviceType == service.TypeLongrun {
		runFile := s.Path + "/run"
		if _, err := os.Stat(runFile); os.IsNotExist(err) {
			s.AddLint("run file for longrun does not exist")
		}
	}

	// check if the up file exists
	if serviceType == "oneshot" {
		upFile := s.Path + "/up"
		if _, err := os.Stat(upFile); os.IsNotExist(err) {
			s.AddLint("up file for oneshot does not exist")
		}
	}

	// check if the dependency directory exists
	dependencyDir, err := s.DependencyDir()
	if err != nil {
		s.AddLint(fmt.Sprintf("service type (%s) in type file for %s does not exist", serviceType, id))
	}

	// get the dependencies
	files, err := os.ReadDir(dependencyDir)
	if err != nil {
		s.AddLint("service has no dependency directory")
		return nil
	}
	for _, file := range files {
		if file.Name() == "base" {
			// we don't want to check the base directory
			continue
		}
		s.AddDependency(file.Name())
	}

	for _, dependency := range s.Dependencies {
		if fs.contains(dependency) || dependency == "base" {
			continue
		}
		// recursive call
		err := fs.compile(dependency)
		if err != nil {
			s.AddLint(fmt.Sprintf("dependency error: %s", err))
		}
	}

	return nil
}

func (fs *Filesystem) serviceType(s *service.Service) (service.Type, error) {
	// check the type file and content
	typeFileContent, err := os.ReadFile(s.Path + "/type")
	if err != nil {
		s.AddLint(fmt.Sprintf("type file for \"%s\" does not exist", s.Id))
		return "", err
	}
	// Check if the file is empty
	if len(typeFileContent) == 0 {
		s.AddLint(fmt.Sprintf("type file for \"%s\" is empty", s.Id))
		return "", err
	}

	// Check if the last character is a newline
	lastChar := typeFileContent[len(typeFileContent)-1]
	if lastChar != '\n' {
		s.AddLint(fmt.Sprintf("type file for \"%s\" does not end with a newline", s.Id))
	}

	serviceType := service.Type(strings.ReplaceAll(string(typeFileContent), "\n", ""))
	for _, theType := range []service.Type{service.TypeOneshot, service.TypeLongrun, service.TypeBundle} {
		if serviceType == theType {
			return serviceType, nil
		}
	}

	return "", errors.New(fmt.Sprintf("invalid type in %s/type file specified", s.Id))
}
