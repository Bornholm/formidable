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

type pointerDeleteTestCase struct {
	DocPath             string
	Pointer             string
	ExpectedRawDocument string
}

func TestPointerDelete(t *testing.T) {
	t.Parallel()

	testCases := []pointerDeleteTestCase{
		{
			DocPath:             "./testdata/set/basic.json",
			Pointer:             "/foo",
			ExpectedRawDocument: `{}`,
		},
		{
			DocPath: "./testdata/set/nested.json",
			Pointer: "/nestedObject/foo/1",
			ExpectedRawDocument: `
			{
				"nestedObject": {
					"foo": [
						"bar",
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
			ExpectedRawDocument: `
			{
			  "nestedObject": {
			    "foo": [
			      "bar",
				  0
			    ]
			  }
			}`,
		},
	}

	for i, tc := range testCases {
		func(index int, tc pointerDeleteTestCase) {
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

				updatedDoc, err := pointer.Delete(baseDoc)
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
					t.Errorf("Delete pointer '%s': expected document \n'%s', got \n'%s'", tc.Pointer, strings.TrimSpace(tc.ExpectedRawDocument), rawDoc)
				}
			})
		}(i, tc)
	}
}
