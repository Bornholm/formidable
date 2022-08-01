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
	return &stdoutFakeCloser{}, nil
}

func NewUpdaterHandler() *UpdaterHandler {
	return &UpdaterHandler{}
}

type stdoutFakeCloser struct {
	io.WriteCloser
}

func (c *stdoutFakeCloser) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (c *stdoutFakeCloser) Close() error {
	return nil
}
