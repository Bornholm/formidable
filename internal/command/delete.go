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

func Delete() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete value at specific path",
		Flags: commonFlags(),
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

			pointer := jsonpointer.New(rawPointer)

			var updatedValues interface{}

			updatedValues, err = pointer.Delete(values)
			if err != nil {
				return errors.Wrapf(err, "could not delete pointer '%v'", rawPointer)
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
