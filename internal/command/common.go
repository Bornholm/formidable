package command

import (
	"bytes"
	"io"
	"net/url"
	"os"
	"reflect"

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
	"forge.cadoles.com/wpetit/formidable/internal/merge"
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
			Usage:   "Use `defaults_url` as defaults",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "values",
			Aliases: []string{"v"},
			Usage:   "Use `values_url` as values",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "schema",
			Aliases: []string{"s"},
			Usage:   "Use `schema_url` as schema",
			Value:   "",
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

func loadData(ctx *cli.Context) (defaults interface{}, values interface{}, err error) {
	values, err = loadValues(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not load values")
	}

	defaults, err = loadDefaults(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not load defaults")
	}

	merged, err := getMatchingZeroValue(values)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	if defaults != nil {
		if err := merge.Merge(&merged, defaults, values); err != nil {
			return nil, nil, errors.Wrap(err, "could not merge values")
		}

		values = merged
	}

	return defaults, values, nil
}

func getMatchingZeroValue(values interface{}) (interface{}, error) {
	valuesKind := reflect.TypeOf(values).Kind()

	switch valuesKind {
	case reflect.Map:
		return make(map[string]interface{}, 0), nil
	case reflect.Slice:
		return make([]interface{}, 0), nil
	default:
		return nil, errors.Errorf("unexpected type '%T'", values)
	}
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
