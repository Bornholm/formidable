package data

import (
	"io"
	"net/url"

	"github.com/pkg/errors"
)

type EncoderHandler interface {
	URLMatcher
	Encode(url *url.URL, data interface{}) (io.Reader, error)
}

type Encoder struct {
	handlers []EncoderHandler
}

func (e *Encoder) Encode(url *url.URL, data interface{}) (io.Reader, error) {
	for _, h := range e.handlers {
		if !h.Match(url) {
			continue
		}

		reader, err := h.Encode(url, data)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return reader, nil
	}

	return nil, errors.Wrapf(ErrHandlerNotFound, "could not find matching handler for url '%s'", url.String())
}

func NewEncoder(handlers ...EncoderHandler) *Encoder {
	return &Encoder{handlers}
}
