package null

import (
	"io"
	"net/url"
)

const SchemeNull = "null"

type UpdaterHandler struct{}

func (h *UpdaterHandler) Match(url *url.URL) bool {
	return url.Scheme == SchemeNull
}

func (u *UpdaterHandler) Update(url *url.URL) (io.WriteCloser, error) {
	return &nullCloser{}, nil
}

func NewUpdaterHandler() *UpdaterHandler {
	return &UpdaterHandler{}
}

type nullCloser struct {
	io.WriteCloser
}

func (c *nullCloser) Write(p []byte) (n int, err error) {
	return io.Discard.Write(p)
}

func (c *nullCloser) Close() error {
	return nil
}
