package internal_test

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/kdisneur/yamlsort/internal"
)

var updateFlag = flag.Bool("update-golden", false, "update golden files")

func TestSorting(t *testing.T) {
	tcs := []struct {
		Name                   string
		Indent                 int
		InputFilePath          string
		ExpectedOutputFilePath string
		Err                    string
	}{
		{
			Name:                   "when YAML is a complete object with bigger indentation",
			Indent:                 4,
			InputFilePath:          "testdata/object-flat-bigger-indent-input.yaml",
			ExpectedOutputFilePath: "testdata/object-flat-bigger-indent-output.yaml",
		},
		{
			Name:                   "when YAML is a complete object",
			Indent:                 2,
			InputFilePath:          "testdata/object-flat-input.yaml",
			ExpectedOutputFilePath: "testdata/object-flat-output.yaml",
		},
		{
			Name:                   "when YAML is part of a bigger object (so indented)",
			Indent:                 2,
			InputFilePath:          "testdata/object-nested-input.yaml",
			ExpectedOutputFilePath: "testdata/object-nested-output.yaml",
		},
		{
			Name:                   "when YAML is a list (so invalid YAML)",
			Indent:                 2,
			InputFilePath:          "testdata/list-input.yaml",
			ExpectedOutputFilePath: "testdata/list-output.yaml",
		},
		{
			Name:                   "when YAML is a list (so invalid YAML) and indented",
			Indent:                 2,
			InputFilePath:          "testdata/list-nested-input.yaml",
			ExpectedOutputFilePath: "testdata/list-nested-output.yaml",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			input, err := os.Open(tc.InputFilePath)
			if err != nil {
				t.Fatalf("can't open input '%s': %s", tc.InputFilePath, err)
			}
			defer input.Close()

			if updateFlag != nil && *updateFlag {
				output, err := os.OpenFile(tc.ExpectedOutputFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
				if err != nil {
					t.Fatalf("can't open golden file '%s' for re-generation: %s", tc.ExpectedOutputFilePath, err)
				}
				defer output.Close()

				if err := internal.SortYAML(input, output, tc.Indent); err != nil {
					t.Fatalf("can't generate golden file: %s", err)
				}
				output.Close()

				if _, err := input.Seek(0, io.SeekStart); err != nil {
					t.Fatalf("can't rewind input file to start after golden file generation: %s", err)
				}
			}

			expectedOutputContent, err := ioutil.ReadFile(tc.ExpectedOutputFilePath)
			if err != nil {
				t.Fatalf("can't read output '%s': %s", tc.ExpectedOutputFilePath, err)
			}

			var actualOutput bytes.Buffer

			err = internal.SortYAML(input, &actualOutput, tc.Indent)
			if err != nil && tc.Err == "" {
				t.Fatalf("expected no error but got: %s", err)
			}

			if err == nil && tc.Err != "" {
				t.Fatalf("expected error '%s' but got none", tc.Err)
			}

			if tc.Err != "" && tc.Err != err.Error() {
				t.Fatalf("wrong error. want: %s; got: %s", tc.Err, err)
			}

			if !bytes.Equal(expectedOutputContent, actualOutput.Bytes()) {
				t.Fatalf("wrong output. want:\n%s\ngot:\n%s", string(expectedOutputContent), actualOutput.String())
			}
		})
	}
}
