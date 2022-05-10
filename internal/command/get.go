package command

import (
	"fmt"
	"os"

	"forge.cadoles.com/wpetit/formidable/internal/jsonpointer"
	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/urfave/cli/v2"
)

func Get() *cli.Command {
	flags := []cli.Flag{}

	flags = append(flags, commonFlags()...)

	return &cli.Command{
		Name:  "get",
		Usage: "Get value at specific path",
		Flags: flags,
		Action: func(ctx *cli.Context) error {
			schema, err := loadSchema(ctx)
			if err != nil {
				return errors.Wrap(err, "could not load schema")
			}

			values, err := loadValues(ctx)
			if err != nil {
				return errors.Wrap(err, "could not load values")
			}

			if err := schema.Validate(values); err != nil {
				if _, ok := err.(*jsonschema.ValidationError); ok {
					fmt.Printf("%#v\n", err)

					os.Exit(1)
				}

				return errors.Wrap(err, "could not validate resulting json")
			}

			rawPointer := ctx.Args().Get(0)
			pointer := jsonpointer.New(rawPointer)

			value, err := pointer.Get(values)
			if err != nil {
				return errors.Wrapf(err, "could not get value from pointer '%v'", rawPointer)
			}

			if err := outputValues(ctx, value); err != nil {
				return errors.Wrap(err, "could not output updated values")
			}

			return nil
		},
	}
}
