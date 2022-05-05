package json

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
)

type parserHandlerTestCase struct {
	Path             string
	ExpectMatch      bool
	ExpectParseError bool
}

var parserHandlerTestCases = []parserHandlerTestCase{
	{
		Path:             "testdata/dummy.json",
		ExpectMatch:      true,
		ExpectParseError: false,
	},
	{
		Path:             "file://testdata/dummy_no_ext?format=json",
		ExpectMatch:      true,
		ExpectParseError: false,
	},
}

func TestDecoderHandler(t *testing.T) {
	t.Parallel()

	handler := NewDecoderHandler()

	for _, tc := range parserHandlerTestCases {
		func(tc parserHandlerTestCase) {
			t.Run(fmt.Sprintf("Parse '%s'", tc.Path), func(t *testing.T) {
				t.Parallel()

				url, err := url.Parse(tc.Path)
				if err != nil {
					t.Fatal(errors.Wrapf(err, "could not parse url '%s'", tc.Path))
				}

				if e, g := tc.ExpectMatch, handler.Match(url); e != g {
					t.Errorf("URL '%s': expected matching result '%v', got '%v'", url.String(), e, g)
				}

				if !tc.ExpectMatch {
					return
				}

				cleanedPath := filepath.Join(url.Host, url.Path)

				file, err := os.Open(cleanedPath)
				if err != nil {
					t.Fatal(errors.Wrapf(err, "could not open file '%s'", cleanedPath))
				}

				defer func() {
					if err := file.Close(); err != nil {
						t.Error(errors.Wrapf(err, "could not close file '%s'", cleanedPath))
					}
				}()

				if _, err := handler.Decode(url, file); err != nil && !tc.ExpectParseError {
					t.Fatal(errors.Wrapf(err, "could not parse file '%s'", tc.Path))
				}

				if tc.ExpectParseError {
					t.Fatal(errors.Errorf("no error was returned as expected when opening url '%s'", url.String()))
				}
			})
		}(tc)
	}
}
