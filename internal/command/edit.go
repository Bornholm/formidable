package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"forge.cadoles.com/wpetit/formidable/internal/server"
	"github.com/pkg/errors"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/urfave/cli/v2"
)

func Edit() *cli.Command {
	flags := commonFlags()

	flags = append(flags, &cli.StringFlag{
		Name:    "browser",
		EnvVars: []string{"FORMIDABLE_BROWSER"},
		Value:   "w3m",
	})

	return &cli.Command{
		Name:  "edit",
		Usage: "Display a form for given schema and values",
		Flags: flags,
		Action: func(ctx *cli.Context) error {
			browser := ctx.String("browser")

			schema, err := loadSchema(ctx)
			if err != nil {
				return errors.Wrap(err, "could not load schema")
			}

			values, err := loadValues(ctx)
			if err != nil {
				return errors.Wrap(err, "could not load values")
			}

			defaults, err := loadDefaults(ctx)
			if err != nil {
				return errors.Wrap(err, "could not load defaults")
			}

			srvCtx, srvCancel := context.WithCancel(ctx.Context)
			defer srvCancel()

			srv := server.New(
				server.WithSchema(schema),
				server.WithValues(values),
				server.WithDefaults(defaults),
			)

			addrs, srvErrs := srv.Start(srvCtx)

			url := fmt.Sprintf("http://%s", (<-addrs).String())
			url = strings.Replace(url, "0.0.0.0", "127.0.0.1", 1)

			log.Printf("listening on %s", url)

			cmdErrs := make(chan error)
			cmdCtx, cmdCancel := context.WithCancel(ctx.Context)
			defer cmdCancel()

			go func() {
				defer func() {
					close(cmdErrs)
				}()

				cmd := exec.CommandContext(cmdCtx, browser, url)

				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				cmd.Env = os.Environ()

				if err := cmd.Run(); err != nil {
					cmdErrs <- errors.WithStack(err)
				}
			}()

			select {
			case err := <-cmdErrs:
				srvCancel()

				return errors.WithStack(err)

			case err := <-srvErrs:
				cmdCancel()

				return errors.WithStack(err)
			}
		},
	}
}
