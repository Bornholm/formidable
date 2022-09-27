package server

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"forge.cadoles.com/wpetit/formidable/internal/jsonpointer"
	"forge.cadoles.com/wpetit/formidable/internal/server/template"
	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func (s *Server) serveFormReq(w http.ResponseWriter, r *http.Request) {
	data := &template.FormItemData{
		Parent:   nil,
		Schema:   s.schema,
		Property: "",
		Defaults: s.defaults,
		Values:   s.values,
	}

	if err := s.schema.Validate(data.Values); err != nil {
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

func (s *Server) handleFormReq(w http.ResponseWriter, r *http.Request) {
	data := &template.FormItemData{
		Parent:   nil,
		Schema:   s.schema,
		Property: "",
		Defaults: s.defaults,
		Values:   s.values,
	}

	var values interface{}

	if err := r.ParseForm(); err != nil {
		panic(errors.WithStack(err))
	} else {
		values, err = handleForm(r.Form, s.schema, s.values)
		if err != nil {
			panic(errors.WithStack(err))
		}

		data.Values = values
	}

	if err := s.schema.Validate(values); err != nil {
		validationErr, ok := err.(*jsonschema.ValidationError)
		if !ok {
			panic(errors.Wrap(err, "could not validate values"))
		}

		data.Error = validationErr
	}

	if data.Error == nil {
		if s.onUpdate != nil {
			if err := s.onUpdate(values); err != nil {
				panic(errors.Wrap(err, "could not update values"))
			}
		}

		data.SuccessMessage = "Data updated."
	}

	if err := template.Exec("index.html.tmpl", w, data); err != nil {
		panic(errors.WithStack(err))
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
			if fieldValues[0] == "" {
				continue
			}

			booVal, err := parseBoolean(fieldValues[0])
			if err != nil {
				return nil, errors.Wrapf(err, "could not parse boolean field '%s'", property)
			}

			pointer := jsonpointer.New(property)

			values, err = pointer.Force(values, booVal)
			if err != nil {
				return nil, errors.Wrapf(err, "could not set property '%s' with value '%v'", property, fieldValues[0])
			}

		case "num":
			if fieldValues[0] == "" {
				continue
			}

			numVal, err := parseNumeric(fieldValues[0])
			if err != nil {
				return nil, errors.Wrapf(err, "could not parse numeric field '%s'", property)
			}

			pointer := jsonpointer.New(property)

			values, err = pointer.Force(values, numVal)
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

func parseNumeric(value string) (float64, error) {
	numVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return numVal, nil
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
