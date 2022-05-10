package stdout

import (
	"io"
	"net/url"
	"os"
)

const SchemeStdout = "stdout"

type UpdaterHandler struct{}

func (h *UpdaterHandler) Match(url *url.URL) bool {
	return url.Scheme == SchemeStdout
}

func (u *UpdaterHandler) Update(url *url.URL) (io.WriteCloser, error) {
	return os.Stdout, nil
}

func NewUpdaterHandler() *UpdaterHandler {
	return &UpdaterHandler{}
}
