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
	Force               bool
	Value               interface{}
	ExpectedError       error
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
						"test",
						{
							"prop1": {
								"subProp": 1
							}
						}
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
						{
							"prop1": {
								"subProp": 1
							}
						},
						"baz"
					]
				}
			}`,
		},
		{
			DocPath: "./testdata/set/nested.json",
			Pointer: "/nestedObject/foo/2/prop2",
			Value:   "baz",
			Force:   true,
			ExpectedRawDocument: `
			{
				"nestedObject": {
					"foo": [
						"bar",
						0,
						{
							"prop2": "baz",
							"prop1": {
								"subProp": 1
							}
						}
					]
				}
			}`,
		},
		{
			DocPath:       "./testdata/set/nested.json",
			Pointer:       "/nestedObject/foo/2/prop2",
			Value:         "baz",
			Force:         false,
			ExpectedError: ErrNotFound,
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

				var updatedDoc interface{}

				if tc.Force {
					updatedDoc, err = pointer.Force(baseDoc, tc.Value)
				} else {
					updatedDoc, err = pointer.Set(baseDoc, tc.Value)
				}

				if tc.ExpectedError != nil && !errors.Is(err, tc.ExpectedError) {
					t.Fatalf("Expected error '%v', got '%v'", tc.ExpectedError, errors.Cause(err))
				}

				if tc.ExpectedError == nil && err != nil {
					t.Fatal(errors.WithStack(err))
				}

				rawDoc, err := json.MarshalIndent(updatedDoc, "", "  ")
				if err != nil {
					t.Fatal(errors.WithStack(err))
				}

				var expectedDoc interface{}

				if tc.ExpectedRawDocument == "" {
					return
				}

				if err := json.Unmarshal([]byte(tc.ExpectedRawDocument), &expectedDoc); err != nil {
					t.Fatal(errors.WithStack(err))
				}

				if !reflect.DeepEqual(expectedDoc, updatedDoc) {
					command := "Set"
					if tc.Force {
						command = "Force"
					}
					t.Errorf("%s pointer '%s' -> '%v': expected document '%s', got '%s'", command, tc.Pointer, tc.Value, strings.TrimSpace(tc.ExpectedRawDocument), rawDoc)
				}
			})
		}(i, tc)
	}
}
