package data

import (
	"io"
	"net/url"

	"github.com/pkg/errors"
)

type UpdaterHandler interface {
	URLMatcher
	Update(url *url.URL) (io.WriteCloser, error)
}

type Updater struct {
	handlers []UpdaterHandler
}

func (u *Updater) Update(url *url.URL) (io.WriteCloser, error) {
	for _, h := range u.handlers {
		if !h.Match(url) {
			continue
		}

		wr, err := h.Update(url)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return wr, nil
	}

	return nil, errors.Wrapf(ErrHandlerNotFound, "could not find matching handler for url '%s'", url.String())
}

func NewUpdater(handlers ...UpdaterHandler) *Updater {
	return &Updater{handlers}
}
