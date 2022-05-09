package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pkg/errors"
)

const (
	testDataDir = "../testdata"
	dummyPath   = "dummy.txt"
)

type loaderHandlerTestCase struct {
	URL               string
	ExpectMatch       bool
	ExpectOpenError   bool
	ExpectOpenContent string
}

func TestLoaderHandler(t *testing.T) {
	t.Parallel()

	staticHandler := http.FileServer(http.Dir(testDataDir))

	server := httptest.NewServer(staticHandler)
	defer server.Close()

	loaderHandlerTestCases := []loaderHandlerTestCase{
		{
			URL:               server.URL + "/" + dummyPath,
			ExpectMatch:       true,
			ExpectOpenError:   false,
			ExpectOpenContent: "dummy",
		},
	}

	handler := NewLoaderHandler(server.Client())

	for _, tc := range loaderHandlerTestCases {
		func(tc loaderHandlerTestCase) {
			t.Run(fmt.Sprintf("Load '%s'", tc.URL), func(t *testing.T) {
				url, err := url.Parse(tc.URL)
				if err != nil {
					t.Fatal(errors.Wrapf(err, "could not parse url '%s'", tc.URL))
				}

				if e, g := tc.ExpectMatch, handler.Match(url); e != g {
					t.Errorf("URL '%s': expected matching result '%v', got '%v'", tc.URL, e, g)
				}

				if !tc.ExpectMatch {
					return
				}

				reader, err := handler.Open(url)
				if err != nil && !tc.ExpectOpenError {
					t.Fatal(errors.Wrapf(err, "could not open url '%s'", url.String()))
				}

				defer func() {
					if err := reader.Close(); err != nil {
						t.Error(errors.WithStack(err))
					}
				}()

				if tc.ExpectOpenError {
					t.Fatal(errors.Errorf("no error was returned as expected when opening url '%s'", url.String()))
				}

				data, err := io.ReadAll(reader)
				if err != nil {
					t.Fatal(errors.WithStack(err))
				}

				if e, g := tc.ExpectOpenContent, string(data); e != g {
					t.Errorf("URL '%s': expected content'%v', got '%v'", tc.URL, e, g)
				}
			})
		}(tc)
	}
}
