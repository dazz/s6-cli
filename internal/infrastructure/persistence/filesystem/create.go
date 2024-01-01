package filesystem

import (
	"fmt"
	"github.com/dazz/s6-cli/internal/domain/service"
	"log"
	"os"
	"strings"
)

func (fs *Filesystem) Oneshot(service *service.Service) error {

	// TODO check if directory exists

	// Create the directory
	if err := os.Mkdir(fs.ServicePath(service.Id), 0755); err != nil && !os.IsExist(err) {
		return err
	}

	// Create the type file
	typeFile := fs.ServicePath(service.Id) + "/type"
	typeFileContent := []byte("oneshot\n")

	if err := os.WriteFile(typeFile, typeFileContent, 0644); err != nil {
		return err
	}

	// Get the absolute path of the specified directory
	scriptFile, err := fs.ServiceScriptFilePath(service.Id)
	if err != nil {
		return err
	}

	// Create the up file
	runFile := fs.ServicePath(service.Id) + "/up"
	runFileContent := []byte(scriptFile + "\n")

	if err := os.WriteFile(runFile, runFileContent, 0644); err != nil {
		return err
	}

	// Create the script file
	scriptFileContent := []byte("#!/bin/sh\n\necho \"hello\"\n")

	if err := os.WriteFile(scriptFile, scriptFileContent, 0755); err != nil {
		return err
	}

	// Create the dependencies
	if err := fs.createDependencies(service); err != nil {
		return err
	}
	return nil
}

// Longrun creates a longrun service
func (fs *Filesystem) Longrun(service *service.Service) error {
	// Create the directory
	if err := os.Mkdir(fs.ServicePath(service.Id), 0755); err != nil {
		return err
	}

	// Create the type file
	typeFile := fs.ServicePath(service.Id) + "/type"
	typeFileContent := []byte("longrun\n")
	if err := os.WriteFile(typeFile, typeFileContent, 0644); err != nil {
		return err
	}

	// Create the run file
	runFile := fs.ServicePath(service.Id) + "/run"
	runFileContent := []byte(fmt.Sprintf("#!/bin/sh\n\nexec %s\n", strings.Join(os.Args[1:], " ")))
	if err := os.WriteFile(runFile, runFileContent, 0644); err != nil {
		return err
	}

	// Create the dependencies
	if err := fs.createDependencies(service); err != nil {
		return err
	}
	return nil
}

// Bundle creates a bundle service
func (fs *Filesystem) Bundle(service *service.Service) error {
	// Create the directory
	if err := os.Mkdir(service.Path, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	// Create the type file
	typeFile := service.Path + "/type"
	typeFileContent := []byte("bundle\n")

	if err := os.WriteFile(typeFile, typeFileContent, 0644); err != nil {
		return err
	}

	// Create the dependencies
	if err := fs.createDependencies(service); err != nil {
		return err
	}
	return nil
}

func (fs *Filesystem) createDependencies(service *service.Service) error {
	dependenciesDir := fs.ServicePath(service.Id) + "/dependencies.d"

	// Create the dependencies directory
	if err := os.Mkdir(dependenciesDir, 0755); err != nil {
		return err
	}
	// Create the dependencies
	for _, dependency := range service.Dependencies {
		// Create the type file
		depFile := dependenciesDir + "/" + string(dependency)
		depFileContent := []byte("")

		if err := os.WriteFile(depFile, depFileContent, 0644); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
