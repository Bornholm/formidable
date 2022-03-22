package command

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"forge.cadoles.com/wpetit/formidable/internal/def"
	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/urfave/cli/v2"
)

const (
	filePathPrefix = "@"
)

func commonFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "defaults",
			Aliases: []string{"d"},
			Usage:   "Default values as JSON or file path prefixed by '@'",
			Value:   "{}",
		},
		&cli.StringFlag{
			Name:    "values",
			Aliases: []string{"v"},
			Usage:   "Current values as JSON or file path prefixed by '@'",
			Value:   "{}",
		},
		&cli.StringFlag{
			Name:      "schema",
			Aliases:   []string{"s"},
			Usage:     "Use `schema_file` as schema",
			TakesFile: true,
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o", "out"},
			Value:   "-",
			Usage:   "Output modified values to `output_file` (or '-' for stdout, the default)",
		},
	}
}

func loadJSONFlag(ctx *cli.Context, flagName string) (interface{}, error) {
	flagValue := ctx.String(flagName)

	if flagValue == "" {
		return nil, nil
	}

	if !strings.HasPrefix(flagValue, filePathPrefix) {
		var value interface{}

		if err := json.Unmarshal([]byte(flagValue), &value); err != nil {
			return nil, errors.WithStack(err)
		}

		return value, nil
	}

	flagValue = strings.TrimPrefix(flagValue, filePathPrefix)

	file, err := os.Open(flagValue)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(errors.WithStack(err))
		}
	}()

	reader := json.NewDecoder(file)

	var values interface{}

	if err := reader.Decode(&values); err != nil {
		return nil, errors.WithStack(err)
	}

	return values, nil
}

func loadValues(ctx *cli.Context) (interface{}, error) {
	values, err := loadJSONFlag(ctx, "values")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return values, nil
}

func loadDefaults(ctx *cli.Context) (interface{}, error) {
	values, err := loadJSONFlag(ctx, "defaults")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return values, nil
}

func loadSchema(ctx *cli.Context) (*jsonschema.Schema, error) {
	schemaFlag := ctx.String("schema")

	compiler := jsonschema.NewCompiler()

	compiler.ExtractAnnotations = true
	compiler.AssertFormat = true
	compiler.AssertContent = true

	var (
		schema *jsonschema.Schema
		err    error
	)

	if schemaFlag == "" {
		schema = def.Schema
	} else {
		schema, err = compiler.Compile(schemaFlag)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return schema, nil
}

const OutputStdout = "-"

type noopWriteCloser struct {
	io.Writer
}

func (c *noopWriteCloser) Close() error {
	return nil
}

func outputWriter(ctx *cli.Context) (io.WriteCloser, error) {
	output := ctx.String("output")

	if output == OutputStdout {
		return &noopWriteCloser{ctx.App.Writer}, nil
	}

	file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return file, nil
}
