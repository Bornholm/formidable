package merge

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type mergeTestCase struct {
	Name           string
	Dst            interface{}
	Dflts          interface{}
	ExpectedResult interface{}
	ShouldFail     bool
}

var mergeTestCases = []mergeTestCase{
	{
		Name: "simple-maps",
		Dst: map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": "test",
			},
		},
		Dflts: map[string]interface{}{
			"other": true,
			"foo": map[string]interface{}{
				"bar": true,
				"baz": 1,
			},
		},
		ExpectedResult: map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": "test",
				"baz": 1,
			},
			"other": true,
		},
	},
	{
		Name:           "string-slices",
		Dst:            []string{"foo"},
		Dflts:          []string{"bar"},
		ExpectedResult: []string{"foo"},
	},
}

func TestMerge(t *testing.T) {
	t.Parallel()

	for _, tc := range mergeTestCases {
		func(tc *mergeTestCase) {
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()

				err := Merge(&tc.Dst, tc.Dflts)

				if tc.ShouldFail && err == nil {
					t.Error("merge should have failed")
				}

				if !reflect.DeepEqual(tc.Dst, tc.ExpectedResult) {
					t.Errorf("tc.Dst should have been the same as tc.ExpectedResult. Expected: %s, got %s", spew.Sdump(tc.ExpectedResult), spew.Sdump(tc.Dst))
				}
			})
		}(&tc)
	}
}
