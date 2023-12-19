package s6cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Service is a struct for a service
type Service struct {
	Name         string
	Type         string
	Dependencies []string
}

// Lint is a struct for a lint
type Lint struct {
	Service string
	Message string
}

func Compile(rootPath string, currentService string, services *[]Service, lints *[]Lint) bool {
	// Check if the service is already in the services list
	if containsService(*services, currentService) {
		return true
	}

	isValid := true

	// Check if the directory exists
	currentServicePath := rootPath + "/" + currentService
	if _, err := os.Stat(currentServicePath); os.IsNotExist(err) {
		*lints = append(*lints, Lint{
			Service: currentService,
			Message: fmt.Sprintf("directory for %s does not exist", currentService),
		})
		return false
	}

	// check the type file and content
	typeFile := currentServicePath + "/type"
	typeFileContent, err := os.ReadFile(typeFile)
	if err != nil {
		*lints = append(*lints, Lint{
			Service: currentService,
			Message: fmt.Sprintf("type file for \"%s\" does not exist", currentService),
		})
		return false
	}
	// Check if the file is empty
	if len(typeFileContent) == 0 {
		*lints = append(*lints, Lint{
			Service: currentService,
			Message: fmt.Sprintf("type file for \"%s\" is empty", currentService),
		})
		return false
	}

	// Check if the last character is a newline
	lastChar := typeFileContent[len(typeFileContent)-1]
	if lastChar != '\n' {
		*lints = append(*lints, Lint{
			Service: currentService,
			Message: fmt.Sprintf("type file for \"%s\" does not end with a newline", currentService),
		})
		isValid = false
	}

	serviceType := strings.ReplaceAll(string(typeFileContent), "\n", "")

	// check if the dependency directory exists
	dependenciesDir := ""
	if serviceType == "bundle" {
		dependenciesDir = currentServicePath + "/contents.d"
	} else if serviceType == "oneshot" || serviceType == "longrun" {
		dependenciesDir = currentServicePath + "/dependencies.d"
	} else {
		*lints = append(*lints, Lint{
			Service: currentService,
			Message: fmt.Sprintf("service type (%s) in type file for %s does not exist", serviceType, currentService),
		})
		return false
	}

	// check if the run file exists
	if serviceType == "longrun" {
		runFile := currentServicePath + "/run"
		if _, err := os.Stat(runFile); os.IsNotExist(err) {
			*lints = append(*lints, Lint{
				Service: currentService,
				Message: fmt.Sprintf("run file for longrun \"%s\" does not exist", currentService),
			})
			isValid = false
		}
	}

	// check if the up file exists
	if serviceType == "oneshot" {
		upFile := currentServicePath + "/up"
		if _, err := os.Stat(upFile); os.IsNotExist(err) {
			*lints = append(*lints, Lint{
				Service: currentService,
				Message: fmt.Sprintf("up file for oneshot \"%s\" does not exist", currentService),
			})
			isValid = false
		}
	}

	// get the dependencies
	files, err := ioutil.ReadDir(dependenciesDir)
	if err != nil {
		log.Fatal(err)
	}
	var dependencies []string
	for _, file := range files {
		dependencies = append(dependencies, file.Name())
		if file.Name() == "base" {
			// we don't want to check the base directory
			continue
		}
	}

	// add the service to the services list
	*services = append(*services, Service{
		Name:         currentService,
		Type:         serviceType,
		Dependencies: dependencies,
	})

	for _, dependency := range dependencies {
		if containsService(*services, dependency) || dependency == "base" {
			continue
		}
		// recursive call
		Compile(rootPath, dependency, services, lints)
	}

	return isValid
}

func MermaidGraph(services []Service) string {
	var graph string
	graph = "```mermaid\ngraph TD;\n"
	for _, service := range services {
		for _, dependency := range service.Dependencies {
			graph += fmt.Sprintf("    %s --> %s\n", service.Name, dependency)
		}
	}
	graph += "```\n"
	return graph
}

func containsService(services []Service, service string) bool {
	for _, s := range services {
		if s.Name == service {
			return true
		}
	}
	return false
}
