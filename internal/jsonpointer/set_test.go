package jsonpointer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

type pointerSetTestCase struct {
	DocPath             string
	Pointer             string
	Value               interface{}
	ExpectedRawDocument string
}

func TestPointerSet(t *testing.T) {
	t.Parallel()

	testCases := []pointerSetTestCase{
		{
			DocPath:             "./testdata/set/basic.json",
			Pointer:             "/foo",
			Value:               "bar",
			ExpectedRawDocument: `{"foo":"bar"}`,
		},
		{
			DocPath: "./testdata/set/nested.json",
			Pointer: "/nestedObject/foo/1",
			Value:   "test",
			ExpectedRawDocument: `
			{
				"nestedObject": {
					"foo": [
						"bar",
						"test"
					]
				}
			}`,
		},
		{
			DocPath: "./testdata/set/nested.json",
			Pointer: "/nestedObject/foo/-",
			Value:   "baz",
			ExpectedRawDocument: `
			{
				"nestedObject": {
					"foo": [
						"bar",
						0,
						"baz"
					]
				}
			}`,
		},
	}

	for i, tc := range testCases {
		func(index int, tc pointerSetTestCase) {
			t.Run(fmt.Sprintf("#%d: '%s'", i, tc.Pointer), func(t *testing.T) {
				t.Parallel()

				baseRawDocument, err := ioutil.ReadFile(tc.DocPath)
				if err != nil {
					t.Fatal(errors.WithStack(err))
				}

				var baseDoc interface{}

				if err := json.Unmarshal([]byte(baseRawDocument), &baseDoc); err != nil {
					t.Fatal(errors.WithStack(err))
				}

				pointer := New(tc.Pointer)

				updatedDoc, err := pointer.Set(baseDoc, tc.Value)
				if err != nil {
					t.Fatal(errors.WithStack(err))
				}

				rawDoc, err := json.MarshalIndent(updatedDoc, "", "  ")
				if err != nil {
					t.Fatal(errors.WithStack(err))
				}

				var expectedDoc interface{}

				if err := json.Unmarshal([]byte(tc.ExpectedRawDocument), &expectedDoc); err != nil {
					t.Fatal(errors.WithStack(err))
				}

				if !reflect.DeepEqual(expectedDoc, updatedDoc) {
					t.Errorf("Set pointer '%s' -> '%v': expected document '%s', got '%s'", tc.Pointer, tc.Value, strings.TrimSpace(tc.ExpectedRawDocument), rawDoc)
				}
			})
		}(i, tc)
	}
}
