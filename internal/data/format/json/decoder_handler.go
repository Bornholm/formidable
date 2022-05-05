package json

import (
	"encoding/json"
	"io"
	"net/url"
	"path"
	"path/filepath"

	"forge.cadoles.com/wpetit/formidable/internal/data/format"
	"github.com/pkg/errors"
)

const (
	ExtensionJSON = ".json"
	FormatJSON    = "json"
)

type DecoderHandler struct{}

func (d *DecoderHandler) Match(url *url.URL) bool {
	ext := filepath.Ext(path.Join(url.Host, url.Path))

	return ext == ExtensionJSON ||
		format.MatchURLQueryFormat(url, FormatJSON)
}

func (d *DecoderHandler) Decode(url *url.URL, reader io.Reader) (interface{}, error) {
	decoder := json.NewDecoder(reader)

	var values interface{}

	if err := decoder.Decode(&values); err != nil {
		return nil, errors.WithStack(err)
	}

	return values, nil
}

func NewDecoderHandler() *DecoderHandler {
	return &DecoderHandler{}
}
