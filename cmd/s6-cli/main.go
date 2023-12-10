package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
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
                Name:    "config",
                Aliases: []string{"c"},
                Value:   "config.yaml",
                Usage:   "Load configuration from `FILE`",
            },
            &cli.StringFlag{
                Name:    "path",
                Aliases: []string{"p"},
                Value:   "/etc/s6-overlay",
                Usage:   "Path to s6-overlay directory",
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

                    // print lints
                    for _, lint := range lints {
                        fmt.Printf("* %s: %s\n", lint.Service, lint.Message)
                    }

                    if valid {
                        return nil
                    }
                    os.Exit(1)
                    return nil
                },
            },
            {
                Name:    "document",
                Aliases: []string{"d"},
                Usage:   "document s6 service dependencies",
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
    // Check if the directory exists
    currentServicePath := rootPath + "/" + currentService
    if _, err := os.Stat(currentServicePath); os.IsNotExist(err) {
        *lints = append(*lints, Lint{
            Service: currentService,
            Message: fmt.Sprintf("directory %s does not exist", currentServicePath),
        })
        return false
    }

    // check the type file and content
    typeFile := currentServicePath + "/type"
    serviceTyp, err := ioutil.ReadFile(typeFile)
    if err != nil {
        *lints = append(*lints, Lint{
            Service: currentService,
            Message: fmt.Sprintf("type file for \"%s\" does not exist", currentService),
        })
        return false
    }
    serviceType := string(serviceTyp)

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
    	// recursive call
    	compileDependencyTree(rootPath, f.Name(), services, lints)
    }

    // add the service to the services list
    *services = append(*services, Service{
        Name: currentService,
        Type: serviceType,
        Dependencies: dependencies,
    })

    return true
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