package json

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"path"
	"path/filepath"

	"forge.cadoles.com/wpetit/formidable/internal/data/format"
	"github.com/pkg/errors"
)

type EncoderHandler struct{}

func (d *EncoderHandler) Match(url *url.URL) bool {
	ext := filepath.Ext(path.Join(url.Host, url.Path))

	return ext == ExtensionJSON ||
		format.MatchURLQueryFormat(url, FormatJSON)
}

func (d *EncoderHandler) Encode(url *url.URL, data interface{}) (io.Reader, error) {
	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)

	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return nil, errors.WithStack(err)
	}

	return &buf, nil
}

func NewEncoderHandler() *EncoderHandler {
	return &EncoderHandler{}
}
