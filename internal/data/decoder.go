package data

import (
	"io"
	"net/url"

	"github.com/pkg/errors"
)

type DecoderHandler interface {
	URLMatcher
	Decode(url *url.URL, reader io.Reader) (interface{}, error)
}

type Decoder struct {
	handlers []DecoderHandler
}

func (d *Decoder) Decode(url *url.URL, reader io.ReadCloser) (interface{}, error) {
	for _, h := range d.handlers {
		if !h.Match(url) {
			continue
		}

		data, err := h.Decode(url, reader)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return data, nil
	}

	return nil, errors.Wrapf(ErrHandlerNotFound, "could not find matching handler for url '%s'", url.String())
}

func NewDecoder(handlers ...DecoderHandler) *Decoder {
	return &Decoder{handlers}
}
