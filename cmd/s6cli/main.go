package main

import (
	"fmt"
	"github.com/dazz/s6-cli/internal/domain/lint"
	"github.com/dazz/s6-cli/internal/infrastructure/persistence"
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
		Name:     "s6-cli",
		Version:  "0.0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Anne-Julia Seitz",
				Email: "dazz@c-base",
			},
		},
		Usage: "CLI for creating and linting files and directories",
		// We'll be using the same flag for all our commands
		// so we'll define it up here
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "rootPath",
				Aliases: []string{"p"},
				Value:   "/etc/s6-overlay/s6-rc.d",
				Usage:   "Path to s6-rc.d directory",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "lint-new",
				Aliases: []string{"ll"},
				Usage:   "lint directories and files",
				Action: func(cCtx *cli.Context) error {
					rootPath := "/etc/s6-overlay/s6-rc.d"

					if cCtx.IsSet("rootPath") {
						rootPath = cCtx.String("rootPath")
					}

					repo := persistence.NewFilesystem(rootPath)
					action := lint.NewAction(repo)

					action.Lint()

					return nil
				},
			},
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
					for _, l := range lints {
						fmt.Printf("* %s: %s\n", l.Service, l.Message)
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
			{
				Name:      "create",
				Aliases:   []string{"c"},
				Usage:     "create a service",
				ArgsUsage: "[type: (o|l|b)] [name]",
				Flags: []cli.Flag{
					&cli.BoolFlag{Value: false, Name: "overwrite", Aliases: []string{"o"}, Usage: "Ignore existing files and directories"},
				},
				Action: func(cCtx *cli.Context) error {
					path := "/etc/s6-overlay/s6-rc.d"
					if cCtx.IsSet("path") {
						path = cCtx.String("path")
					}
					// check if the directory exists
					if _, err := os.Stat(path); os.IsNotExist(err) {
						fmt.Printf("Directory %s does not exist\n", path)
						os.Exit(1)
					}

					serviceType := ""
					if cCtx.Args().Get(0) == "o" || cCtx.Args().Get(0) == "l" || cCtx.Args().Get(0) == "b" {
						serviceType = cCtx.Args().Get(0)
					} else {
						fmt.Print("Arg type must not be empty and one of 'o', 'l' or 'b'\n")
						os.Exit(1)
					}

					name := ""
					if cCtx.Args().Get(1) != "" {
						name = cCtx.Args().Get(1)
					} else {
						fmt.Print("Arg name must not be empty\n")
						os.Exit(1)
					}

					fmt.Printf("Create a service %s with type %s\n", name, serviceType)

					if serviceType == "o" {
						s6cli.Oneshot(path, name, []string{"base"})
					}

					return nil
				},
			},
			{
				Name:      "remove",
				Aliases:   []string{"rm"},
				Usage:     "remove a service",
				ArgsUsage: "[name]",
				Action: func(cCtx *cli.Context) error {
					path := "/etc/s6-overlay/s6-rc.d"
					if cCtx.IsSet("path") {
						path = cCtx.String("path")
					}
					// check if the directory exists
					if _, err := os.Stat(path); os.IsNotExist(err) {
						fmt.Printf("Directory %s does not exist\n", path)
						os.Exit(1)
					}

					name := ""
					if cCtx.Args().Get(0) != "" {
						name = cCtx.Args().Get(0)
					} else {
						fmt.Print("Arg name must not be empty\n")
						os.Exit(1)
					}

					s6cli.Remove(path, name)

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
