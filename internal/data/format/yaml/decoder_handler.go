package yaml

import (
	"io"
	"net/url"
	"path"
	"path/filepath"
	"regexp"

	"forge.cadoles.com/wpetit/formidable/internal/data/format"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

var (
	ExtensionYAML = regexp.MustCompile("\\.ya?ml$")
	FormatYAML    = "yaml"
)

type DecoderHandler struct{}

func (d *DecoderHandler) Match(url *url.URL) bool {
	ext := filepath.Ext(path.Join(url.Host, url.Path))

	return ExtensionYAML.MatchString(ext) ||
		format.MatchURLQueryFormat(url, FormatYAML)
}

func (d *DecoderHandler) Decode(url *url.URL, reader io.Reader) (interface{}, error) {
	decoder := yaml.NewDecoder(reader)

	var values interface{}

	if err := decoder.Decode(&values); err != nil {
		return nil, errors.WithStack(err)
	}

	return values, nil
}

func NewDecoderHandler() *DecoderHandler {
	return &DecoderHandler{}
}
