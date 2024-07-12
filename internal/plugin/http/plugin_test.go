package http_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/555f/gg/pkg/gg"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testData struct {
	String string
}

func Test(t *testing.T) {
	wd, _ := filepath.Abs("./tests")

	type output struct {
		path         string
		expectedFile string
	}

	testCases := []struct {
		name     string
		fileName string
		outputs  []output
	}{
		{
			name:     "Echo no wrapper error",
			fileName: "echo_no_wrapper_error.go",
			outputs: []output{
				{
					"/home/vitaly/Documents/work/my/gg/internal/plugin/http/tests/internal/server/client.go",
					"echo_no_wrapper_error.go",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results, err := gg.Run("dev", wd, []string{"./" + tc.fileName}, map[string]string{}, false)
			if err != nil {
				t.Error(err)
			}
			if len(results) != len(tc.outputs) {
				t.Error("results not equal outputs")
			}
			for i, r := range results {
				if actual, expected := r.File.Path(), tc.outputs[i].path; actual != expected {
					t.Errorf("path actual %s\nexpected %s", actual, expected)
				}
				actualBytes, err := r.File.Bytes()
				if err != nil {
					t.Error(err)
				}
				expectedFiepath := filepath.Join(wd, "testdata", tc.outputs[i].expectedFile)
				if _, err := os.Stat(expectedFiepath); err != nil {
					f, err := os.Create(filepath.Join(wd, "testdata", tc.fileName))
					if err != nil {
						t.Error(err)
					}
					_, err = f.Write(actualBytes)
					if err != nil {
						t.Error(err)
					}
					f.Close()
					t.Skip()
				}
				expectedBytes, err := os.ReadFile(expectedFiepath)
				if err != nil {
					t.Error(err)
				}

				actual, expected := string(actualBytes), string(expectedBytes)
				diff := cmp.Diff(actual, expected, cmpopts.AcyclicTransformer("multiline", func(s string) []string {
					return strings.Split(s, "\n")
				}))
				if diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
