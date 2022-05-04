package template

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path"
	"strings"

	"forge.cadoles.com/wpetit/formidable/internal/jsonpointer"
	"github.com/Masterminds/sprig/v3"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

var (
	layouts map[string]*template.Template
	blocks  map[string]string
)

func Load(files embed.FS, baseDir string) error {
	if blocks == nil {
		blocks = make(map[string]string)
	}

	blockFiles, err := fs.ReadDir(files, path.Join(baseDir, "blocks"))
	if err != nil {
		return errors.WithStack(err)
	}

	for _, f := range blockFiles {
		templateData, err := fs.ReadFile(files, path.Join(baseDir, "blocks/"+f.Name()))
		if err != nil {
			return errors.WithStack(err)
		}

		blocks[f.Name()] = string(templateData)
	}

	layoutFiles, err := fs.ReadDir(files, path.Join(baseDir, "layouts"))
	if err != nil {
		return errors.WithStack(err)
	}

	for _, f := range layoutFiles {
		templateData, err := fs.ReadFile(files, path.Join(baseDir, "layouts/"+f.Name()))
		if err != nil {
			return errors.WithStack(err)
		}

		if err := loadLayout(f.Name(), string(templateData)); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func loadLayout(name string, rawTemplate string) error {
	if layouts == nil {
		layouts = make(map[string]*template.Template)
	}

	tmpl := template.New(name)
	funcMap := mergeHelpers(
		sprig.FuncMap(),
		customHelpers(tmpl),
	)

	tmpl.Funcs(funcMap)

	for blockName, b := range blocks {
		if _, err := tmpl.Parse(b); err != nil {
			return errors.Wrapf(err, "could not parse template block '%s'", blockName)
		}
	}

	tmpl, err := tmpl.Parse(rawTemplate)
	if err != nil {
		return errors.Wrapf(err, "could not parse template '%s'", name)
	}

	layouts[name] = tmpl

	return nil
}

func Exec(name string, w io.Writer, data interface{}) error {
	tmpl, exists := layouts[name]
	if !exists {
		return errors.Errorf("could not find template '%s'", name)
	}

	if err := tmpl.Execute(w, data); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func mergeHelpers(helpers ...template.FuncMap) template.FuncMap {
	merged := template.FuncMap{}

	for _, help := range helpers {
		for name, fn := range help {
			merged[name] = fn
		}
	}

	return merged
}

type FormItemData struct {
	Parent         *FormItemData
	Schema         *jsonschema.Schema
	Property       string
	Error          *jsonschema.ValidationError
	Values         interface{}
	Defaults       interface{}
	SuccessMessage string
}

func customHelpers(tpl *template.Template) template.FuncMap {
	return template.FuncMap{
		"formItemData": func(parent *FormItemData, property string, schema *jsonschema.Schema) *FormItemData {
			return &FormItemData{
				Parent:   parent,
				Property: property,
				Schema:   schema,
				Defaults: parent.Defaults,
				Values:   parent.Values,
				Error:    parent.Error,
			}
		},
		"dump": func(data interface{}) string {
			spew.Dump(data)

			return ""
		},
		"include": func(name string, data interface{}) (template.HTML, error) {
			buf := bytes.NewBuffer([]byte{})

			if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
				return "", errors.WithStack(err)
			}

			return template.HTML(buf.String()), nil
		},
		"getFullProperty": func(parent *FormItemData, property string) string {
			fullProperty := property
			for {
				fullProperty = fmt.Sprintf("%s/%s", parent.Property, strings.TrimPrefix(fullProperty, "/"))
				parent = parent.Parent
				if parent == nil {
					break
				}
			}

			return fullProperty
		},
		"getValue": func(defaults, values interface{}, path string) (interface{}, error) {
			if defaults == nil {
				defaults = make(map[string]interface{})
			}

			if values == nil {
				values = make(map[string]interface{})
			}

			pointer := jsonpointer.New(path)

			val, err := pointer.Get(values)
			if err != nil && !errors.Is(err, jsonpointer.ErrNotFound) {
				return nil, errors.WithStack(err)
			}

			if errors.Is(err, jsonpointer.ErrNotFound) {
				val, err = pointer.Get(defaults)
				if err != nil && !errors.Is(err, jsonpointer.ErrNotFound) {
					return nil, errors.WithStack(err)
				}
			}

			return val, nil
		},
		"getItemSchema": func(arraySchema *jsonschema.Schema) (*jsonschema.Schema, error) {
			itemSchema := arraySchema.Items
			if itemSchema == nil {
				itemSchema = arraySchema.Items2020
			}

			if itemSchema == nil {
				return nil, errors.New("item schema not found")
			}

			switch schema := itemSchema.(type) {
			case *jsonschema.Schema:
				return schema, nil
			case []*jsonschema.Schema:
				if len(schema) > 0 {
					return schema[0], nil
				}

				return nil, errors.New("no item schema found")
			default:
				return nil, errors.Errorf("unexpected schema type '%T'", schema)
			}
		},
		"getPropertyError": findPropertyValidationError,
	}
}

func findPropertyValidationError(err *jsonschema.ValidationError, property string) *jsonschema.ValidationError {
	if err == nil {
		return nil
	}

	if property == err.InstanceLocation {
		return err
	}

	for _, cause := range err.Causes {
		if err := findPropertyValidationError(cause, property); err != nil {
			return err
		}
	}

	return nil
}
