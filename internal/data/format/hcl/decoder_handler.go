package hcl

import (
	"io"
	"net/url"
	"path"
	"path/filepath"

	"forge.cadoles.com/wpetit/formidable/internal/data/format"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

const (
	ExtensionHCL = ".hcl"
	FormatHCL    = "hcl"
)

type DecoderHandler struct {
	ctx *hcl.EvalContext
}

func (d *DecoderHandler) Match(url *url.URL) bool {
	ext := filepath.Ext(path.Join(url.Host, url.Path))

	return ext == ExtensionHCL ||
		format.MatchURLQueryFormat(url, FormatHCL)
}

func (d *DecoderHandler) Decode(url *url.URL, reader io.Reader) (interface{}, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ctx := d.ctx
	if ctx == nil {
		ctx = &hcl.EvalContext{
			Variables: make(map[string]cty.Value),
			Functions: make(map[string]function.Function),
		}
	}

	parser := hclparse.NewParser()

	file, diags := parser.ParseHCL(data, url.String())
	if diags.HasErrors() {
		return nil, errors.WithStack(diags)
	}

	var tree map[string]interface{}

	diags = gohcl.DecodeBody(file.Body, ctx, &tree)
	if diags.HasErrors() {
		return nil, errors.WithStack(diags)
	}

	ctx = ctx.NewChild()

	values, err := hclTreeToRawValues(ctx, tree)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return values, nil
}

func NewDecoderHandler(ctx *hcl.EvalContext) *DecoderHandler {
	return &DecoderHandler{ctx}
}

func hclTreeToRawValues(ctx *hcl.EvalContext, tree map[string]interface{}) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	for key, branch := range tree {
		v, err := hclBranchToRawValue(ctx, branch)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		values[key] = v
	}

	return values, nil
}

func hclBranchToRawValue(ctx *hcl.EvalContext, branch interface{}) (interface{}, error) {
	switch typ := branch.(type) {
	case *hcl.Attribute:
		val, diags := typ.Expr.Value(ctx)
		if diags.HasErrors() {
			return nil, errors.WithStack(diags)
		}

		raw, err := ctyValueToRaw(ctx, val)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return raw, nil
	default:
		return nil, errors.Errorf("unexpected type '%T'", typ)
	}
}

func ctyValueToRaw(ctx *hcl.EvalContext, val cty.Value) (interface{}, error) {
	if val.Type().Equals(cty.Bool) {
		var raw bool

		if err := gocty.FromCtyValue(val, &raw); err != nil {
			return nil, errors.WithStack(err)
		}

		return raw, nil
	} else if val.Type().Equals(cty.Number) {
		var raw float64

		if err := gocty.FromCtyValue(val, &raw); err != nil {
			return nil, errors.WithStack(err)
		}

		return raw, nil
	} else if val.Type().Equals(cty.String) {
		var raw string

		if err := gocty.FromCtyValue(val, &raw); err != nil {
			return nil, errors.WithStack(err)
		}

		return raw, nil
	} else if val.Type().IsObjectType() {
		obj := make(map[string]interface{})

		for k, v := range val.AsValueMap() {
			rv, err := ctyValueToRaw(ctx, v)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			obj[k] = rv
		}

		return obj, nil
	} else if val.Type().IsTupleType() {
		sl := make([]interface{}, 0)

		for _, v := range val.AsValueSlice() {
			rv, err := ctyValueToRaw(ctx, v)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			sl = append(sl, rv)
		}

		return sl, nil
	}

	return nil, errors.Errorf("unexpected cty.Type '%s'", val.Type().FriendlyName())
}
