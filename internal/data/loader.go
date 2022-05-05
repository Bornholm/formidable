package data

import (
	"io"
	"net/url"

	"github.com/pkg/errors"
)

type LoaderHandler interface {
	URLMatcher
	Open(url *url.URL) (io.ReadCloser, error)
}

type Loader struct {
	handlers []LoaderHandler
}

func (l *Loader) Open(url *url.URL) (io.ReadCloser, error) {
	for _, h := range l.handlers {
		if !h.Match(url) {
			continue
		}

		reader, err := h.Open(url)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return reader, nil
	}

	return nil, errors.Wrapf(ErrHandlerNotFound, "could not find matching handler for url '%s'", url.String())
}

func NewLoader(handlers ...LoaderHandler) *Loader {
	return &Loader{handlers}
}
