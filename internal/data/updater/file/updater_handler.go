package file

import (
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const SchemeFile = "file"

type UpdaterHandler struct{}

func (h *UpdaterHandler) Match(url *url.URL) bool {
	return url.Scheme == SchemeFile
}

func (u *UpdaterHandler) Update(url *url.URL) (io.WriteCloser, error) {
	name := filepath.Join(url.Host, url.Path)

	file, err := os.Create(name)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open file '%s'", name)
	}

	return file, nil
}

func NewUpdaterHandler() *UpdaterHandler {
	return &UpdaterHandler{}
}
