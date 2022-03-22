package server

import (
	"forge.cadoles.com/wpetit/formidable/internal/def"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

type Option struct {
	Host     string
	Port     uint
	Schema   *jsonschema.Schema
	Values   interface{}
	Defaults interface{}
}

type OptionFunc func(*Option)

func defaultOption() *Option {
	return &Option{
		Host:   "",
		Port:   0,
		Schema: def.Schema,
	}
}

func WithAddress(host string, port uint) OptionFunc {
	return func(opt *Option) {
		opt.Host = host
		opt.Port = port
	}
}

func WithSchema(schema *jsonschema.Schema) OptionFunc {
	return func(opt *Option) {
		opt.Schema = schema
	}
}

func WithValues(values interface{}) OptionFunc {
	return func(opt *Option) {
		opt.Values = values
	}
}

func WithDefaults(defaults interface{}) OptionFunc {
	return func(opt *Option) {
		opt.Defaults = defaults
	}
}
