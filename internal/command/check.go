package command

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/urfave/cli/v2"
)

func Check() *cli.Command {
	flags := []cli.Flag{}

	flags = append(flags, commonFlags()...)

	return &cli.Command{
		Name:  "check",
		Usage: "Check values with the given schema",
		Flags: flags,
		Action: func(ctx *cli.Context) error {
			schema, err := loadSchema(ctx)
			if err != nil {
				return errors.Wrap(err, "could not load schema")
			}

			_, values, err := loadData(ctx)
			if err != nil {
				return errors.Wrap(err, "could not load data")
			}

			if err := schema.Validate(values); err != nil {
				if _, ok := err.(*jsonschema.ValidationError); ok {
					fmt.Printf("%#v\n", err)

					return errors.New("invalid values")
				}

				return errors.Wrap(err, "could not validate values")
			}

			if err := outputValues(ctx, values); err != nil {
				return errors.Wrap(err, "could not output updated values")
			}

			return nil
		},
	}
}
