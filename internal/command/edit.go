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

	flags = append(flags,
		&cli.StringFlag{
			Name:    "browser",
			EnvVars: []string{"FORMIDABLE_BROWSER"},
			Value:   "w3m",
		},
		&cli.StringFlag{
			Name:    "http-host",
			EnvVars: []string{"FORMIDABLE_HTTP_HOST"},
			Value:   "127.0.0.1",
		},
		&cli.UintFlag{
			Name:    "http-port",
			EnvVars: []string{"FORMIDABLE_HTTP_PORT"},
			Value:   0,
		},
		&cli.BoolFlag{
			Name:    "no-browser",
			EnvVars: []string{"FORMIDABLE_NO_BROWSER"},
			Value:   false,
		},
	)

	return &cli.Command{
		Name:  "edit",
		Usage: "Display a form for given schema and values",
		Flags: flags,
		Action: func(ctx *cli.Context) error {
			browser := ctx.String("browser")
			noBrowser := ctx.Bool("no-browser")
			httpPort := ctx.Uint("http-port")
			httpHost := ctx.String("http-host")

			schema, err := loadSchema(ctx)
			if err != nil {
				return errors.Wrap(err, "could not load schema")
			}

			defaults, values, err := loadData(ctx)
			if err != nil {
				return errors.Wrap(err, "could not load data")
			}

			srvCtx, srvCancel := context.WithCancel(ctx.Context)
			defer srvCancel()

			srv := server.New(
				server.WithAddress(httpHost, httpPort),
				server.WithSchema(schema),
				server.WithValues(values),
				server.WithDefaults(defaults),
				server.WithOnUpdate(func(values interface{}) error {
					if err := outputValues(ctx, values); err != nil {
						return errors.Wrap(err, "could not output updated values")
					}

					return nil
				}),
			)

			addrs, srvErrs := srv.Start(srvCtx)

			url := fmt.Sprintf("http://%s", (<-addrs).String())
			url = strings.Replace(url, "0.0.0.0", "127.0.0.1", 1)

			log.Printf("listening on %s", url)

			browserErrs := make(chan error)
			browserCtx, browserCancel := context.WithCancel(ctx.Context)
			defer browserCancel()

			if !noBrowser {
				browserErrs = startBrowser(browserCtx, browser, url)
			}

			select {
			case err := <-browserErrs:
				srvCancel()

				return errors.WithStack(err)

			case err := <-srvErrs:
				browserCancel()

				return errors.WithStack(err)
			}
		},
	}
}

func startBrowser(ctx context.Context, browser, url string) chan error {
	cmdErrs := make(chan error)

	go func() {
		defer func() {
			close(cmdErrs)
		}()

		cmd := exec.CommandContext(ctx, browser, url)

		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Env = os.Environ()

		if err := cmd.Run(); err != nil {
			cmdErrs <- errors.WithStack(err)
		}
	}()

	return cmdErrs
}
