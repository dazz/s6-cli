package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "strings"
    "time"
    "github.com/urfave/cli/v2"
)

func init() {
}

func main() {
    app := &cli.App{
        Name:  "s6-cli",
        Version: "0.0.1",
        Compiled: time.Now(),
        Authors: []*cli.Author{
            &cli.Author{
                Name:  "Anne-Julia Seitz",
                Email: "dazz@c-base",
            },
        },
        Usage: "CLI for creating and linting files and directories",
        // We'll be using the same flag for all our commands
        // so we'll define it up here
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:    "path",
                Aliases: []string{"p"},
                Value:   "/etc/s6-overlay/s6-rc.d",
                Usage:   "Path to s6-rc.d directory",
            },
        },
        Commands: []*cli.Command{
            {
                Name:    "lint",
                Aliases: []string{"l"},
                Usage:   "lint directories and files",
                Action: func(cCtx *cli.Context) error {
                    path := "/etc/s6-overlay/s6-rc.d"
                    firstBundle := "user"

                    if cCtx.IsSet("path") {
                        path = cCtx.String("path")
                    }
                    // check if the directory exists
                    if _, err := os.Stat(path); os.IsNotExist(err) {
                        fmt.Printf("Directory %s does not exist\n", path)
                        os.Exit(1)
                    }

                    // compile dependency tree
                    var services []Service
                    var lints []Lint
                    valid := compileDependencyTree(path, firstBundle, &services, &lints)

                    fmt.Println("*************** s6-cli Lint Report ***************")

                    // print lints
                    for _, lint := range lints {
                        fmt.Printf("* %s: %s\n", lint.Service, lint.Message)
                    }

                    fmt.Println("*************** s6-cli Lint Report ***************")

                    if valid {
                        return nil
                    }
                    os.Exit(1)
                    return nil
                },
            },
            {
                Name:    "mermaid",
                Aliases: []string{"m"},
                Usage:   "document s6 service dependencies in mermaid syntax",
                Action: func(cCtx *cli.Context) error {
                    path := "/etc/s6-overlay/s6-rc.d"
                    firstBundle := "user"

                    if cCtx.IsSet("path") {
                        path = cCtx.String("path")
                    }
                    // check if the directory exists
                    if _, err := os.Stat(path); os.IsNotExist(err) {
                        fmt.Printf("Directory %s does not exist\n", path)
                        os.Exit(1)
                    }

                    // compile dependency tree
                    var services []Service
                    var lints []Lint
                    compileDependencyTree(path, firstBundle, &services, &lints)

                    fmt.Printf(renderMermaidGraph(services))

                    return nil
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

// Define a struct named Person
type Service struct {
    Name string
    Type string
    Dependencies []string
}
type Lint struct {
    Service string
    Message string
}


func compileDependencyTree(rootPath string, currentService string, services *[]Service, lints *[]Lint) bool {
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

    if currentService == "user" {}

    // check the type file and content
    typeFile := currentServicePath + "/type"
    typeFileContent, err := ioutil.ReadFile(typeFile)
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
    for _, f := range files {
    	dependencies = append(dependencies, f.Name())
    	if f.Name() == "base" {
    	    // we don't want to check the base directory
            continue
        }
    }

    // add the service to the services list
    *services = append(*services, Service{
        Name: currentService,
        Type: serviceType,
        Dependencies: dependencies,
    })

    for _, dependency := range dependencies {
        if containsService(*services, dependency) || dependency == "base" {
            continue
        }
        // recursive call
        compileDependencyTree(rootPath, dependency, services, lints)
    }

    return isValid
}

func renderMermaidGraph(services []Service) string {
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

func containsService (services []Service, service string) bool {
    for _, s := range services {
        if s.Name == service {
            return true
        }
    }
    return false
}