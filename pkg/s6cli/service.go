package s6cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Oneshot(rootPath string, name string, dependencies []string) {
	servicePath := rootPath + "/" + name

	// Create the directory
	err := os.Mkdir(servicePath, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// Create the type file
	typeFile := servicePath + "/type"
	typeFileContent := []byte("oneshot\n")
	err = os.WriteFile(typeFile, typeFileContent, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Get the absolute path of the specified directory
	scriptFile := getScriptsPath(rootPath) + "/" + name

	// Create the up file
	runFile := servicePath + "/up"
	runFileContent := []byte(scriptFile + "\n")
	err = os.WriteFile(runFile, runFileContent, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Create the script file
	scriptFileContent := []byte("#!/bin/sh\n\necho \"hello\"\n")
	err = os.WriteFile(scriptFile, scriptFileContent, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Create the dependencies
	dependenciesDir := servicePath + "/dependencies.d"
	createDependencies(dependenciesDir, dependencies)
}

// Longrun creates a longrun service
func Longrun(rootPath string, name string, dependencies []string) {
	servicePath := rootPath + "/" + name

	// Create the directory
	os.Mkdir(servicePath, 0755)

	// Create the type file
	typeFile := servicePath + "/type"
	typeFileContent := []byte("longrun\n")
	os.WriteFile(typeFile, typeFileContent, 0644)

	// Create the run file
	runFile := servicePath + "/run"
	runFileContent := []byte(fmt.Sprintf("#!/bin/sh\n\nexec %s\n", strings.Join(os.Args[1:], " ")))
	os.WriteFile(runFile, runFileContent, 0644)

	// Create the dependencies
	dependenciesDir := servicePath + "/dependencies.d"
	createDependencies(dependenciesDir, dependencies)
}

// Bundle creates a bundle service
func Bundle(rootPath string, name string, dependencies []string) {
	servicePath := rootPath + "/" + name

	// Create the directory
	err := os.Mkdir(servicePath, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// Create the type file
	typeFile := servicePath + "/type"
	typeFileContent := []byte("bundle\n")
	err = os.WriteFile(typeFile, typeFileContent, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Create the dependencies
	dependenciesDir := servicePath + "/contents.d"
	createDependencies(dependenciesDir, dependencies)
}

// Remove removes a service
func Remove(rootPath string, name string) {
	servicePath := rootPath + "/" + name

	if serviceType := getServiceType(servicePath); serviceType == "oneshot" {
		// Remove the script file
		scriptFile := getScriptsPath(rootPath) + "/" + name
		err := os.Remove(scriptFile)
		if err != nil {
			fmt.Println("Error removing script file:", err)
		}
	}

	// Remove the service directory
	err := os.RemoveAll(servicePath)
	if err != nil {
		log.Fatal(err)
	}

	// Remove service from dependency directories
	// Find files with the specified name in the directory and its children
	dependencyFoundAt, err := findFilesByName(rootPath, name)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(dependencyFoundAt) > 0 {
		for _, depFile := range dependencyFoundAt {
			err := os.Remove(depFile)
			if err != nil {
				fmt.Println("Error removing dependency file:", err)
			}
		}
	}
}

func createDependencies(dependenciesDir string, dependencies []string) {
	// Create the dependencies directory
	err := os.Mkdir(dependenciesDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	// Create the dependencies
	for _, dependency := range dependencies {
		// Create the type file
		depFile := dependenciesDir + "/" + dependency
		depFileContent := []byte("")
		err = os.WriteFile(depFile, depFileContent, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getScriptsPath(rootPath string) string {
	// Get the absolute path of the specified directory
	absoluteScriptPath, err := filepath.Abs(rootPath + "/../scripts")
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return ""
	}
	return absoluteScriptPath
}

func getServiceType(servicePath string) string {
	typeFile := servicePath + "/type"
	typeFileContent, err := os.ReadFile(typeFile)
	if err != nil {
		fmt.Println("Error reading type file:", err)
		return ""
	}
	return strings.ReplaceAll(string(typeFileContent), "\n", "")
}

func findFilesByName(root, targetFileName string) ([]string, error) {
	var resultPaths []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current item is a regular file and matches the target file name
		if !info.IsDir() && info.Name() == targetFileName {
			resultPaths = append(resultPaths, path)
		}

		return nil
	})

	return resultPaths, err
}
