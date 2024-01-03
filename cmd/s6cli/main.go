package main

import (
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
				Name:    "lint",
				Aliases: []string{"l"},
				Usage:   "lint directories and files",
				Action: func(cCtx *cli.Context) error {
					if cCtx.IsSet("rootPath") {
						rootPath = cCtx.String("rootPath")
					}

					repo := filesystem.NewFilesystem(rootPath)
					command := lint.NewCommand(repo)

					execute, err := command.Execute()
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(execute)

					return nil
				},
			},
			{
				Name:    "mermaid",
				Aliases: []string{"m"},
				Usage:   "document s6 service dependencies in mermaid syntax",
				Action: func(cCtx *cli.Context) error {
					if cCtx.IsSet("rootPath") {
						rootPath = cCtx.String("rootPath")
					}

					repo := filesystem.NewFilesystem(rootPath)
					command := mermaid.NewCommand(repo)

					execute, err := command.Execute()
					if err != nil {

					}

					fmt.Println(execute)

					return nil
				},
			},
			{
				Name:      "create",
				Aliases:   []string{"c"},
				Usage:     "create a service",
				ArgsUsage: "[type: (o|l|b)] [id]",
				Flags: []cli.Flag{
					&cli.BoolFlag{Value: false, Name: "overwrite", Aliases: []string{"o"}, Usage: "Ignore existing files and directories"},
				},
				Action: func(cCtx *cli.Context) error {
					rootPath := "/etc/s6-overlay/s6-rc.d"
					if cCtx.IsSet("rootPath") {
						rootPath = cCtx.String("rootPath")
					}
					// check if the directory exists
					if _, err := os.Stat(rootPath); os.IsNotExist(err) {
						fmt.Printf("Directory %s does not exist\n", rootPath)
						os.Exit(1)
					}

					var serviceType service.Type
					switch t := cCtx.Args().Get(0); t {
					case "o":
						serviceType = service.TypeOneshot
					case "l":
						serviceType = service.TypeLongrun
					case "b":
						serviceType = service.TypeBundle
					default:
						fmt.Print("Arg type must not be empty and one of 'o', 'l' or 'b'\n")
						os.Exit(1)
					}

					var id service.Id
					if idArg := cCtx.Args().Get(1); idArg != "" {
						id = service.Id(idArg)
					} else {
						fmt.Println("Arg idArg must not be empty")
						os.Exit(1)
					}

					fmt.Printf("Creating a service %q with type %s\n", id, serviceType)

					repo := filesystem.NewFilesystem(rootPath)
					command := create.NewCommand(repo, id, serviceType)

					result, err := command.Execute()
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					fmt.Printf("Successful created service %q\n", result)
					return nil

					return nil
				},
			},
			{
				Name:      "remove",
				Aliases:   []string{"rm"},
				Usage:     "remove a service",
				ArgsUsage: "[name]",
				Action: func(cCtx *cli.Context) error {
					rootPath := "/etc/s6-overlay/s6-rc.d"
					if cCtx.IsSet("rootPath") {
						rootPath = cCtx.String("rootPath")
					}
					// check if the directory exists
					if _, err := os.Stat(rootPath); os.IsNotExist(err) {
						fmt.Println("Directory does not exist: " + rootPath)
						os.Exit(1)
					}

					var id service.Id
					if idArg := cCtx.Args().Get(0); idArg != "" {
						id = service.Id(idArg)
					} else {
						fmt.Println("Arg idArg must not be empty")
						os.Exit(1)
					}

					repo := filesystem.NewFilesystem(rootPath)
					command := remove.NewCommand(repo, id)

					result, err := command.Execute()
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
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
