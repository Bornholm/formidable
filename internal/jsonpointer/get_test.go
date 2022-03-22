package jsonpointer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

type pointerGetTestCase struct {
	Document         interface{}
	Pointer          string
	ExpectedRawValue string
}

func TestPointerGet(t *testing.T) {
	t.Parallel()

	ietfRawDocument, err := ioutil.ReadFile("./testdata/ietf.json")
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	var ietfDoc interface{}

	if err := json.Unmarshal([]byte(ietfRawDocument), &ietfDoc); err != nil {
		t.Fatal(errors.WithStack(err))
	}

	// IETF tests cases
	// From https://datatracker.ietf.org/doc/html/rfc6901
	//
	// ""           // the whole document
	// "/foo"       ["bar", "baz"]
	// "/foo/0"     "bar"
	// "/"          0
	// "/a~1b"      1
	// "/c%d"       2
	// "/e^f"       3
	// "/g|h"       4
	// "/i\\j"      5
	// "/k\"l"      6
	// "/ "         7
	// "/m~0n"      8
	testCases := []pointerGetTestCase{
		{
			Document:         ietfDoc,
			Pointer:          "",
			ExpectedRawValue: string(ietfRawDocument),
		},
		{
			Document:         ietfDoc,
			Pointer:          "/foo",
			ExpectedRawValue: "[\"bar\", \"baz\"]",
		},
		{
			Document:         ietfDoc,
			Pointer:          "/foo/0",
			ExpectedRawValue: "\"bar\"",
		},
		{
			Document:         ietfDoc,
			Pointer:          `/`,
			ExpectedRawValue: `0`,
		},
		{
			Document:         ietfDoc,
			Pointer:          "/a~1b",
			ExpectedRawValue: "1",
		},
		{
			Document:         ietfDoc,
			Pointer:          "/c%d",
			ExpectedRawValue: "2",
		},
		{
			Document:         ietfDoc,
			Pointer:          "/e^f",
			ExpectedRawValue: "3",
		},
		{
			Document:         ietfDoc,
			Pointer:          "/g|h",
			ExpectedRawValue: "4",
		},
		{
			Document:         ietfDoc,
			Pointer:          "/i\\j",
			ExpectedRawValue: "5",
		},
		{
			Document:         ietfDoc,
			Pointer:          "/k\"l",
			ExpectedRawValue: "6",
		},
		{
			Document:         ietfDoc,
			Pointer:          "/ ",
			ExpectedRawValue: "7",
		},
		{
			Document:         ietfDoc,
			Pointer:          "/m~0n",
			ExpectedRawValue: "8",
		},
	}

	for i, tc := range testCases {
		func(index int, tc pointerGetTestCase) {
			t.Run(fmt.Sprintf("#%d: '%s'", i, tc.Pointer), func(t *testing.T) {
				t.Parallel()

				pointer := New(tc.Pointer)

				value, err := pointer.Get(tc.Document)
				if err != nil {
					t.Fatal(errors.WithStack(err))
				}

				rawValue, err := json.Marshal(value)
				if err != nil {
					t.Fatal(errors.WithStack(err))
				}

				var expectedValue interface{}

				if err := json.Unmarshal([]byte(tc.ExpectedRawValue), &expectedValue); err != nil {
					t.Fatal(errors.WithStack(err))
				}

				if !reflect.DeepEqual(expectedValue, value) {
					t.Errorf("Pointer '%s': expected value '%s', got '%s'", tc.Pointer, tc.ExpectedRawValue, rawValue)
				}
			})
		}(i, tc)
	}
}
