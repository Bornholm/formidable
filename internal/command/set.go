package command

import (
	"encoding/json"
	"fmt"
	"os"

	"forge.cadoles.com/wpetit/formidable/internal/jsonpointer"
	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/urfave/cli/v2"
)

func Set() *cli.Command {
	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "Force data tree creation",
			Value:   false,
		},
	}

	flags = append(flags, commonFlags()...)

	return &cli.Command{
		Name:  "set",
		Usage: "Set value at specific path",
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

			rawPointer := ctx.Args().Get(0)
			rawValue := ctx.Args().Get(1)

			pointer := jsonpointer.New(rawPointer)

			var value interface{}

			if err := json.Unmarshal([]byte(rawValue), &value); err != nil {
				return errors.Wrapf(err, "could not parse json '%s'", rawValue)
			}

			var updatedValues interface{}

			force := ctx.Bool("force")

			if force {
				updatedValues, err = pointer.Force(values, value)
				if err != nil {
					return errors.Wrapf(err, "could not force value '%v' to pointer '%v'", rawValue, rawPointer)
				}
			} else {
				updatedValues, err = pointer.Set(values, value)
				if err != nil {
					return errors.Wrapf(err, "could not set value '%v' to pointer '%v'", rawValue, rawPointer)
				}
			}

			if err := schema.Validate(updatedValues); err != nil {
				if _, ok := err.(*jsonschema.ValidationError); ok {
					fmt.Printf("%#v\n", err)

					os.Exit(1)
				}

				return errors.Wrap(err, "could not validate resulting json")
			}

			output, err := outputWriter(ctx)
			if err != nil {
				return errors.Wrap(err, "could not create output writer")
			}

			encoder := json.NewEncoder(output)

			encoder.SetIndent("", "  ")

			if err := encoder.Encode(updatedValues); err != nil {
				return errors.Wrap(err, "could not write to output")
			}

			return nil
		},
	}
}
