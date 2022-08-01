package command

import (
	"bytes"
	"io"
	"net/url"
	"os"

	encjson "encoding/json"

	"forge.cadoles.com/wpetit/formidable/internal/data"
	"forge.cadoles.com/wpetit/formidable/internal/data/format/hcl"
	"forge.cadoles.com/wpetit/formidable/internal/data/format/json"
	"forge.cadoles.com/wpetit/formidable/internal/data/format/yaml"
	"forge.cadoles.com/wpetit/formidable/internal/data/scheme/file"
	"forge.cadoles.com/wpetit/formidable/internal/data/scheme/http"
	"forge.cadoles.com/wpetit/formidable/internal/data/scheme/stdin"
	"forge.cadoles.com/wpetit/formidable/internal/data/updater/exec"
	fileUpdater "forge.cadoles.com/wpetit/formidable/internal/data/updater/file"
	"forge.cadoles.com/wpetit/formidable/internal/data/updater/null"
	"forge.cadoles.com/wpetit/formidable/internal/data/updater/stdout"
	"forge.cadoles.com/wpetit/formidable/internal/def"
	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/urfave/cli/v2"

	gohttp "net/http"
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
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "values",
			Aliases: []string{"v"},
			Usage:   "Current values as JSON or file path prefixed by '@'",
			Value:   "",
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
			Value:   "stdout://local?format=json",
			Usage:   "Output modified values to specified URL",
		},
	}
}

func loadURLFlag(ctx *cli.Context, flagName string) (interface{}, error) {
	flagValue := ctx.String(flagName)

	if flagValue == "" {
		return nil, nil
	}

	loader := newLoader()

	url, err := url.Parse(flagValue)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	reader, err := loader.Open(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		if err := reader.Close(); err != nil {
			panic(errors.WithStack(err))
		}
	}()

	decoder := newDecoder()

	data, err := decoder.Decode(url, reader)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return data, nil
}

func loadValues(ctx *cli.Context) (interface{}, error) {
	values, err := loadURLFlag(ctx, "values")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return values, nil
}

func loadDefaults(ctx *cli.Context) (interface{}, error) {
	defaults, err := loadURLFlag(ctx, "defaults")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return defaults, nil
}

func loadSchema(ctx *cli.Context) (*jsonschema.Schema, error) {
	schemaFlag := ctx.String("schema")

	if schemaFlag == "" {
		return def.Schema, nil
	}

	schemaTree, err := loadURLFlag(ctx, "schema")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Reencode schema to JSON format
	var buf bytes.Buffer
	encoder := encjson.NewEncoder(&buf)

	if err := encoder.Encode(schemaTree); err != nil {
		return nil, errors.WithStack(err)
	}

	compiler := jsonschema.NewCompiler()

	compiler.ExtractAnnotations = true
	compiler.AssertFormat = true
	compiler.AssertContent = true

	if err := compiler.AddResource(schemaFlag, &buf); err != nil {
		return nil, errors.WithStack(err)
	}

	schema, err := compiler.Compile(schemaFlag)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return schema, nil
}

func outputValues(ctx *cli.Context, values interface{}) error {
	outputFlag := ctx.String("output")

	url, err := url.Parse(outputFlag)
	if err != nil {
		return errors.WithStack(err)
	}

	encoder := newEncoder()

	reader, err := encoder.Encode(url, values)
	if err != nil {
		return errors.WithStack(err)
	}

	updater := newUpdater()

	writer, err := updater.Update(url)
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		if err := writer.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
			panic(errors.WithStack(err))
		}
	}()

	if _, err := io.Copy(writer, reader); err != nil && !errors.Is(err, io.EOF) {
		return errors.WithStack(err)
	}

	return nil
}

func newLoader() *data.Loader {
	return data.NewLoader(
		file.NewLoaderHandler(),
		http.NewLoaderHandler(gohttp.DefaultClient),
		stdin.NewLoaderHandler(),
	)
}

func newDecoder() *data.Decoder {
	return data.NewDecoder(
		json.NewDecoderHandler(),
		hcl.NewDecoderHandler(nil),
		yaml.NewDecoderHandler(),
	)
}

func newUpdater() *data.Updater {
	return data.NewUpdater(
		stdout.NewUpdaterHandler(),
		fileUpdater.NewUpdaterHandler(),
		exec.NewUpdaterHandler(),
		null.NewUpdaterHandler(),
	)
}

func newEncoder() *data.Encoder {
	return data.NewEncoder(
		json.NewEncoderHandler(),
		yaml.NewEncoderHandler(),
	)
}
