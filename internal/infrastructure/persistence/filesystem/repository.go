package filesystem

import (
	"errors"
	"fmt"
	"github.com/dazz/s6-cli/internal/domain/service"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const RootService service.Id = "user"

type Filesystem struct {
	rootPath string
	byId     map[service.Id]*service.Service
	allIds   []service.Id
}

func NewFilesystem(rootPath string) *Filesystem {

	// check if the directory exists
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		log.Fatal("rootPath directory does not exist: " + rootPath)
	}

	return &Filesystem{
		rootPath: rootPath,
		byId:     make(map[service.Id]*service.Service),
		allIds:   []service.Id{},
	}
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

	s := service.NewService(id)

	// add the service to the services list
	fs.byId[id] = s
	fs.allIds = append(fs.allIds, id)

	// Check if the directory exists
	if !fs.FileExists(fs.ServicePath(id)) {
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
		runFile := fs.ServicePath(id) + "/run"
		if _, err := os.Stat(runFile); os.IsNotExist(err) {
			s.AddLint("run file for longrun does not exist")
		}
	}

	// check if the up file exists
	if serviceType == "oneshot" {
		upFile := fs.ServicePath(id) + "/up"
		if _, err := os.Stat(upFile); os.IsNotExist(err) {
			s.AddLint("up file for oneshot does not exist")
		}
	}

	// check if the dependency directory exists
	dependencyDir, err := fs.ServiceDependenciesPath(s)
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
	typeFileContent, err := os.ReadFile(fs.ServicePath(s.Id) + "/type")
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

	return "", fmt.Errorf("invalid type in %s/type file specified", s.Id)
}

func (fs *Filesystem) FindDependenciesById(id service.Id) ([]string, error) {
	var resultPaths []string

	err := filepath.Walk(fs.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current item is a regular file and matches the target file name
		if !info.IsDir() && info.Name() == string(id) {
			resultPaths = append(resultPaths, path)
		}

		return nil
	})

	return resultPaths, err
}

func (fs *Filesystem) FileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func (fs *Filesystem) AllFileExists(files []string) bool {
	for _, file := range files {
		if !fs.FileExists(file) {
			return false
		}
	}
	return true
}

func (fs *Filesystem) AnyFileExists(files []string) bool {
	for _, file := range files {
		if fs.FileExists(file) {
			return true
		}
	}
	return false
}

func (fs *Filesystem) ServiceScriptFilePath(id service.Id) (string, error) {
	// Get the absolute path of the specified directory
	absoluteScriptPath, err := fs.ScriptPath()
	if err != nil {
		return "", fmt.Errorf("error getting absolute path: %s", err)
	}
	return fmt.Sprintf("%s/%s", absoluteScriptPath, id), nil
}

func (fs *Filesystem) ScriptPath() (string, error) {
	absoluteScriptPath, err := filepath.Abs(fs.rootPath + "/../scripts")
	return absoluteScriptPath, err
}

func (fs *Filesystem) ServicePath(id service.Id) string {
	if id == "" {
		return fs.rootPath
	}
	return fs.rootPath + "/" + string(id)
}

func (fs *Filesystem) ServiceDependenciesPath(s *service.Service) (string, error) {
	if s.Type == "" {
		return "", errors.New("invalid service type, set type of service")
	}

	switch s.Type {
	case service.TypeLongrun:
		return fs.ServicePath(s.Id) + "/dependencies.d", nil
	case service.TypeOneshot:
		return fs.ServicePath(s.Id) + "/dependencies.d", nil
	case service.TypeBundle:
		return fs.ServicePath(s.Id) + "/contents.d", nil
	}
	return "", errors.New("invalid service type")
}
