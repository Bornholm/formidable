package yaml

import (
	"bytes"
	"io"
	"net/url"
	"path"
	"path/filepath"

	"forge.cadoles.com/wpetit/formidable/internal/data/format"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type EncoderHandler struct{}

func (d *EncoderHandler) Match(url *url.URL) bool {
	ext := filepath.Ext(path.Join(url.Host, url.Path))

	return ExtensionYAML.MatchString(ext) ||
		format.MatchURLQueryFormat(url, FormatYAML)
}

func (d *EncoderHandler) Encode(url *url.URL, data interface{}) (io.Reader, error) {
	var buf bytes.Buffer

	encoder := yaml.NewEncoder(&buf)

	if err := encoder.Encode(data); err != nil {
		return nil, errors.WithStack(err)
	}

	return &buf, nil
}

func NewEncoderHandler() *EncoderHandler {
	return &EncoderHandler{}
}
