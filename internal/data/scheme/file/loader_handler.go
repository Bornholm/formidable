package file

import (
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const SchemeFile = "file"

type LoaderHandler struct{}

func (h *LoaderHandler) Match(url *url.URL) bool {
	return url.Scheme == SchemeFile || url.Scheme == ""
}

func (h *LoaderHandler) Open(url *url.URL) (io.ReadCloser, error) {
	file, err := os.Open(filepath.Join(url.Host, url.Path))
	if err != nil {
		return nil, errors.Wrapf(err, "could not open file '%s'", url.Path)
	}

	return file, nil
}

func NewLoaderHandler() *LoaderHandler {
	return &LoaderHandler{}
}
