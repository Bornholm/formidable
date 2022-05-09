package stdin

import (
	"io"
	"net/url"
	"os"
)

const SchemeStdin = "stdin"

type LoaderHandler struct{}

func (h *LoaderHandler) Match(url *url.URL) bool {
	return url.Scheme == SchemeStdin
}

func (h *LoaderHandler) Open(url *url.URL) (io.ReadCloser, error) {
	return os.Stdin, nil
}

func NewLoaderHandler() *LoaderHandler {
	return &LoaderHandler{}
}
