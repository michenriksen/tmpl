package testutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const updateGoldenEnv = "UPDATE_GOLDEN"

const (
	kindGolden = "golden"
	kindActual = "actual"
)

// gitDiffArgs are the arguments to git diff when generating a diff between
// golden and actual files.
var gitDiffArgs = []string{
	"--no-pager", "diff", "--no-index", "--color=always", "--src-prefix=./",
	"--dst-prefix=./",
}

// Golden provides functionality for testing with golden files.
//
// Read more on golden file testing: https://ieftimov.com/posts/testing-in-go-golden-files/
type Golden struct {
	tb       testing.TB
	isUpdate bool
}

// NewGolden returns a new golden file test helper for the given testing.TB.
//
// Golden files are stored in the testdata/golden directory. Each test has a
// golden file with the same name as the test function. The golden file contains
// the expected output of the test. When the test is run, the actual output is
// compared to the golden file. If the actual output does not match the golden
// file, the actual output is written to a file in the same directory as the
// golden file for debugging purposes.
func NewGolden(tb testing.TB) *Golden {
	tb.Helper()

	_, update := os.LookupEnv(updateGoldenEnv)

	return &Golden{tb: tb, isUpdate: update}
}

// AssertMatch asserts that the provided value matches the expected, golden
// value.
//
// If the value does not match, the test will fail and the actual value will be
// written to a file in the same directory as the golden file for debugging
// purposes. If the values differ because of an expected change, the file can
// be renamed to the golden file, or the UPDATE_GOLDEN environment variable can
// be set to update the golden file.
//
// NOTE: the golden and actual values are compared by their marshaled JSON form
// and are also written to files in this format. This makes it easy to read and
// compare the file contents.
func (g *Golden) AssertMatch(got any) {
	g.tb.Helper()

	actual := g.marshal(got)
	expected := g.readOrWriteGolden(actual)

	// If the test has already failed, don't bother comparing the values.
	if g.tb.Failed() {
		return
	}

	if !bytes.Equal(actual, expected) {
		g.writeActual(actual)

		g.tb.Errorf("expected test data for %s to match golden data:\n\n%s\n", g.tb.Name(), g.diff())
	}
}

// RequireMatch is like [AssertMatch], but will fail the test immediately if the
// provided value does not match the expected, golden value.
func (g *Golden) RequireMatch(got any) {
	g.tb.Helper()

	if g.AssertMatch(got); g.tb.Failed() {
		g.tb.FailNow()
	}
}

// readOrWriteGolden reads the golden file if it exists, otherwise it writes the
// actual value to the golden file if the UPDATE_GOLDEN environment variable is
// set.
func (g *Golden) readOrWriteGolden(actual []byte) []byte {
	g.tb.Helper()

	gPath := g.path(kindGolden)

	info, err := os.Stat(gPath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			g.tb.Fatalf("getting file info on golden file %q: %v", gPath, err)
		}

		if !g.isUpdate {
			g.tb.Errorf("golden file for %s does not exist; run with %s=1 to create it", g.tb.Name(), updateGoldenEnv)
			return nil
		}
	}

	if g.isUpdate {
		g.update(actual)
		return actual
	}

	golden, err := os.ReadFile(gPath)
	if err != nil {
		g.tb.Fatalf("reading golden file %q: %v", gPath, err)
	}

	g.tb.Logf("read golden file %q (%dB)", gPath, info.Size())

	return golden
}

func (g *Golden) writeGolden(data []byte) {
	g.tb.Helper()
	WriteFile(g.tb, data, g.path(kindGolden))
}

func (g *Golden) writeActual(data []byte) {
	g.tb.Helper()
	WriteFile(g.tb, data, g.path(kindActual))
}

func (g *Golden) rmActual() {
	g.tb.Helper()

	aPath := g.path(kindActual)
	if err := os.Remove(aPath); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			g.tb.Fatalf("removing actual file %q: %v", aPath, err)
		}

		return
	}

	g.tb.Logf("removed actual file %q", aPath)
}

// update updates the golden file with the provided data and removes the
// '*.actual.json' file if it exists.
func (g *Golden) update(data []byte) {
	g.tb.Helper()

	g.writeGolden(data)
	g.rmActual()
}

func (g *Golden) path(kind string) string {
	g.tb.Helper()
	return filepath.Join("testdata", kindGolden, fmt.Sprintf("%s.%s.json", g.tb.Name(), kind))
}

func (g *Golden) marshal(v any) []byte {
	g.tb.Helper()

	if b, ok := v.([]byte); ok {
		v = g.comparableBytes(b)
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		g.tb.Fatalf("marshaling golden data: %v", err)
	}

	// Add a newline to the end to ensure a blank newline when writing to a file
	// and to make the diff output more readable.
	return append(b, '\n')
}

// comparableBytes transforms the provided byte slice into a value that is more
// easily comparable.
func (g *Golden) comparableBytes(data []byte) any {
	g.tb.Helper()

	ct := http.DetectContentType(data)

	if ct == "application/json" {
		return data
	}

	if strings.HasPrefix(ct, "text/") {
		lines := bytes.Split(data, []byte("\n"))
		linesStr := make([]string, len(lines))

		for i, line := range lines {
			linesStr[i] = string(line)
		}

		return linesStr
	}

	return data
}

func (g *Golden) diff() string {
	g.tb.Helper()

	git, err := exec.LookPath("git")
	if err != nil {
		g.tb.Fatal("git executable is required to diff golden files")
	}

	cmd := exec.Command(git, append(gitDiffArgs, g.path(kindGolden), g.path(kindActual))...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		exitErr := &exec.ExitError{}
		if !errors.As(err, &exitErr) {
			g.tb.Logf("command output: %s", out)
			g.tb.Fatalf("running command %q: %v", cmd.String(), err)
		}
	}

	return string(out)
}
