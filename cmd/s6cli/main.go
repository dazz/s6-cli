package main

import (
	"errors"
	"fmt"
	"github.com/dazz/s6-cli/internal/domain/create"
	"github.com/dazz/s6-cli/internal/domain/lint"
	"github.com/dazz/s6-cli/internal/domain/mermaid"
	"github.com/dazz/s6-cli/internal/domain/remove"
	"github.com/dazz/s6-cli/internal/domain/service"
	"github.com/dazz/s6-cli/internal/infrastructure/persistence/filesystem"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	rootPath := "/etc/s6-overlay/s6-rc.d"
	app := &cli.App{
		Name:     "s6-cli",
		Version:  "v0.2.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Anne-Julia Seitz",
				Email: "dazz@c-base",
			},
		},
		Usage: "CLI for creating and linting files and directories",
		// We'll be using the same flag for all our commands,
		// so we'll define it up here
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "root-path",
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
					if cCtx.IsSet("root-path") {
						rootPath = cCtx.String("root-path")
					}

					repo := filesystem.NewFilesystem(rootPath)
					command := lint.NewCommand(repo)

					execute, err := command.Execute()
					if err != nil {
						return err
					}
					if execute != "" {
						fmt.Println("s6-cli: lint found issues with services in " + rootPath)
						fmt.Println(execute)
						return nil
					}
					fmt.Println("s6-cli: lint found no issues")
					return nil
				},
			},
			{
				Name:    "mermaid",
				Aliases: []string{"m"},
				Usage:   "document s6 service dependencies in mermaid syntax",
				Action: func(cCtx *cli.Context) error {
					if cCtx.IsSet("root-path") {
						rootPath = cCtx.String("root-path")
					}

					repo := filesystem.NewFilesystem(rootPath)
					command := mermaid.NewCommand(repo)

					execute, err := command.Execute()
					if err != nil {
						return err
					}

					fmt.Println(execute)
					return nil
				},
			},
			{
				Name:      "create",
				Aliases:   []string{"c"},
				Usage:     "create a service",
				ArgsUsage: "<type: oneshot|longrun|bundle> <id>",
				Flags: []cli.Flag{
					&cli.BoolFlag{Value: false, Name: "overwrite", Aliases: []string{"o"}, Usage: "Ignore existing files and directories"},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.IsSet("root-path") {
						rootPath = cCtx.String("root-path")
					}

					var serviceType = service.Type(cCtx.Args().Get(0))
					if service.ValidType(serviceType) == false {
						fmt.Println("Argument type be one of 'oneshot', 'longrun' or 'bundle'")
						os.Exit(1)
					}

					var id service.Id
					if idArg := cCtx.Args().Get(1); idArg != "" {
						id = service.Id(idArg)
					} else {
						fmt.Println("Argument id must not be empty")
						os.Exit(1)
					}

					fmt.Printf("Creating a service %q with type %s\n", id, serviceType)

					repo := filesystem.NewFilesystem(rootPath)
					command := create.NewCommand(repo, id, serviceType)

					result, err := command.Execute()
					if err != nil {
						return err
					}
					fmt.Printf("Successful created service %q\n", result)
					return nil
				},
			},
			{
				Name:      "remove",
				Aliases:   []string{"rm"},
				Usage:     "remove a service",
				ArgsUsage: "[name]",
				Action: func(cCtx *cli.Context) error {
					if cCtx.IsSet("root-path") {
						rootPath = cCtx.String("root-path")
					}

					var id service.Id
					if idArg := cCtx.Args().Get(0); idArg != "" {
						id = service.Id(idArg)
					} else {
						return errors.New("arg idArg must not be empty")
					}

					repo := filesystem.NewFilesystem(rootPath)
					command := remove.NewCommand(repo, id)

					result, err := command.Execute()
					if err != nil {
						return err
					}
					fmt.Printf("Successful removed service %q\n", result)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
