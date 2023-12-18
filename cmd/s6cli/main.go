package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/urfave/cli/v2"

    "github.com/dazz/s6-cli/pkg/s6cli"
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
                    var services []s6cli.Service
                    var lints []s6cli.Lint
                    valid := s6cli.Compile(path, firstBundle, &services, &lints)

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
                    var services []s6cli.Service
                    var lints []s6cli.Lint
                    s6cli.Compile(path, firstBundle, &services, &lints)

                    fmt.Printf(s6cli.MermaidGraph(services))

                    return nil
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
