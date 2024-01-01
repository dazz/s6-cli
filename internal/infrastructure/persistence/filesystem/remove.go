package filesystem

import (
	"errors"
	"fmt"
	"github.com/dazz/s6-cli/internal/domain/service"
	"os"
)

// Remove removes a service
func (fs *Filesystem) Remove(s *service.Service) error {

	if s.Type == "oneshot" {
		// Remove the script file
		scriptFilePath, err := fs.ServiceScriptFilePath(s.Id)
		if err != nil {
			return err
		}
		if err := os.Remove(scriptFilePath); err != nil {
			return errors.New(fmt.Sprintf("Error removing script file: %s", err))
		}
	}

	// Remove the service directory
	if err := os.RemoveAll(fs.ServicePath(s.Id)); err != nil {
		return err
	}

	// Remove service from dependency directories
	// Find files with the specified name in the directory and its children
	dependencyFoundAt, err := fs.findDependenciesById(s.Id)
	if err != nil {
		return err
	}

	if len(dependencyFoundAt) > 0 {
		for _, depFile := range dependencyFoundAt {
			if err := os.Remove(depFile); err != nil {
				return errors.New(fmt.Sprintf("Error removing dependency file: %s", err))
			}
		}
	}
	return nil
}
