package http

import (
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	SchemeHTTP  = "http"
	SchemeHTTPS = "https"
)

type LoaderHandler struct {
	client *http.Client
}

func (h *LoaderHandler) Match(url *url.URL) bool {
	return url.Scheme == SchemeHTTP || url.Scheme == SchemeHTTPS
}

func (h *LoaderHandler) Open(url *url.URL) (io.ReadCloser, error) {
	res, err := h.client.Get(url.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code '%d (%s)'", res.StatusCode, http.StatusText(res.StatusCode))
	}

	return res.Body, nil
}

func NewLoaderHandler(client *http.Client) *LoaderHandler {
	return &LoaderHandler{client}
}
