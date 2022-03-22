package route

import (
	"net/http"
	"net/url"
	"strings"

	"forge.cadoles.com/wpetit/formidable/internal/jsonpointer"
	"forge.cadoles.com/wpetit/formidable/internal/server/template"
	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func createRenderFormHandlerFunc(schema *jsonschema.Schema, defaults, values interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &template.FormItemData{
			Parent:   nil,
			Schema:   schema,
			Property: "",
			Defaults: defaults,
			Values:   values,
		}

		if err := schema.Validate(data.Values); err != nil {
			validationErr, ok := err.(*jsonschema.ValidationError)
			if !ok {
				panic(errors.Wrap(err, "could not validate values"))
			}

			data.Error = validationErr
		}

		if err := template.Exec("index.html.tmpl", w, data); err != nil {
			panic(errors.WithStack(err))
		}
	}
}

func createHandleFormHandlerFunc(schema *jsonschema.Schema, defaults, values interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &template.FormItemData{
			Parent:   nil,
			Schema:   schema,
			Property: "",
			Defaults: defaults,
			Values:   values,
		}

		if err := r.ParseForm(); err != nil {
			panic(errors.WithStack(err))
		} else {
			values, err = handleForm(r.Form, schema, values)
			if err != nil {
				panic(errors.WithStack(err))
			}

			data.Values = values
		}

		if err := schema.Validate(data.Values); err != nil {
			validationErr, ok := err.(*jsonschema.ValidationError)
			if !ok {
				panic(errors.Wrap(err, "could not validate values"))
			}

			data.Error = validationErr
		}

		if err := template.Exec("index.html.tmpl", w, data); err != nil {
			panic(errors.WithStack(err))
		}
	}
}

func handleForm(form url.Values, schema *jsonschema.Schema, values interface{}) (interface{}, error) {
	pendingDeletes := make([]string, 0)

	var err error

	for name, fieldValues := range form {
		if name == "submit" {
			continue
		}

		prefix, property, err := parseFieldName(name)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		switch prefix {
		case "bool":
			booVal, err := parseBoolean(fieldValues[0])
			if err != nil {
				return nil, errors.Wrapf(err, "could not parse boolean field '%s'", property)
			}

			pointer := jsonpointer.New(property)

			values, err = pointer.Force(values, booVal)
			if err != nil {
				return nil, errors.Wrapf(err, "could not set property '%s' with value '%v'", property, fieldValues[0])
			}

		case "add":
			pointer := jsonpointer.New(property)

			values, err = pointer.Force(values, nil)
			if err != nil {
				return nil, errors.Wrapf(err, "could not add item '%s'", property)
			}

		case "del":
			// Mark property for deletion pass
			pendingDeletes = append(pendingDeletes, property)

		default:
			pointer := jsonpointer.New(property)

			values, err = pointer.Force(values, fieldValues[0])
			if err != nil {
				return nil, errors.Wrapf(err, "could not set property '%s' with value '%v'", property, fieldValues[0])
			}
		}
	}

	for _, property := range pendingDeletes {
		pointer := jsonpointer.New(property)

		values, err = pointer.Delete(values)
		if err != nil {
			return nil, errors.Wrapf(err, "could not delete property '%s'", property)
		}
	}

	return values, nil
}

func parseBoolean(value string) (bool, error) {
	switch value {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	default:
		return false, errors.Errorf("unexpected boolean value '%s'", value)
	}
}

func parseFieldName(name string) (string, string, error) {
	tokens := strings.SplitN(name, ":", 2)

	if len(tokens) == 1 {
		return "", tokens[0], nil
	}

	if len(tokens) == 2 {
		return tokens[0], tokens[1], nil
	}

	return "", "", errors.Errorf("unexpected field name '%s'", name)
}
