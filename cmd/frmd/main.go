package main

import (
	"context"
	"fmt"
	"os"
	"sort"

	"forge.cadoles.com/wpetit/formidable/internal/command"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// nolint: gochecknoglobals
var (
	GitRef         = "unknown"
	ProjectVersion = "unknown"
	BuildDate      = "unknown"
)

func main() {
	ctx := context.Background()

	app := &cli.App{
		Name:     "frmd",
		Usage:    "JSON Schema based cli forms",
		Commands: command.Root(),
		Before: func(ctx *cli.Context) error {
			workdir := ctx.String("workdir")
			// Switch to new working directory if defined
			if workdir != "" {
				if err := os.Chdir(workdir); err != nil {
					return errors.Wrap(err, "could not change working directory")
				}
			}

			if err := ctx.Set("projectVersion", ProjectVersion); err != nil {
				return errors.WithStack(err)
			}

			if err := ctx.Set("gitRef", GitRef); err != nil {
				return errors.WithStack(err)
			}

			if err := ctx.Set("buildDate", BuildDate); err != nil {
				return errors.WithStack(err)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "workdir",
				Value: "",
				Usage: "The working directory",
			},
			&cli.StringFlag{
				Name:   "projectVersion",
				Value:  "",
				Hidden: true,
			},
			&cli.StringFlag{
				Name:   "gitRef",
				Value:  "",
				Hidden: true,
			},
			&cli.StringFlag{
				Name:   "buildDate",
				Value:  "",
				Hidden: true,
			},
		},
	}

	app.ExitErrHandler = func(ctx *cli.Context, err error) {
		fmt.Printf("%+v", err)
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		panic(errors.WithStack(err))
	}
}
