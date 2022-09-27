package command

import (
	"flag"
	"testing"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type ExpectFunc func(t *testing.T, cmd *cli.Command, err error)

type checkCommandTestCase struct {
	Name        string
	SchemaFile  string
	DefaultFile string
	ValuesFile  string
	Expect      ExpectFunc
}

var checkCommandTestCases = []checkCommandTestCase{
	{
		Name:        "ok",
		SchemaFile:  "file://testdata/check/schema.json",
		DefaultFile: "file://testdata/check/defaults.json",
		ValuesFile:  "file://testdata/check/values-ok.json",
		Expect:      expectNoError,
	},
	{
		Name:        "nok",
		SchemaFile:  "file://testdata/check/schema.json",
		DefaultFile: "file://testdata/check/defaults.json",
		ValuesFile:  "file://testdata/check/values-nok.json",
		Expect:      expectError,
	},
}

func TestCheck(t *testing.T) {
	t.Parallel()

	for _, tc := range checkCommandTestCases {
		func(tc *checkCommandTestCase) {
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()

				flags := flag.NewFlagSet("", flag.ExitOnError)
				cmd := Check()

				for _, f := range cmd.Flags {
					if err := f.Apply(flags); err != nil {
						t.Fatal(errors.WithStack(err))
					}
				}

				err := flags.Parse([]string{
					"check",
					"--schema", tc.SchemaFile,
					"--defaults", tc.DefaultFile,
					"--values", tc.ValuesFile,
					"--output", "null://local?format=json",
				})
				if err != nil {
					t.Fatal(errors.WithStack(err))
				}

				app := cli.NewApp()
				ctx := cli.NewContext(app, flags, nil)

				err = cmd.Run(ctx)

				tc.Expect(t, cmd, err)
			})
		}(&tc)
	}
}

func expectNoError(t *testing.T, cmd *cli.Command, err error) {
	if err != nil {
		t.Error(errors.Wrap(err, "the command result in an unexpected error"))
	}
}

func expectError(t *testing.T, cmd *cli.Command, err error) {
	if err == nil {
		t.Error(errors.New("an error should have been returned"))
	}
}
