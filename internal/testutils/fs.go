package testutils

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func ReadFile(tb testing.TB, pathElems ...string) []byte {
	tb.Helper()

	name := filepath.Join(pathElems...)

	data, err := os.ReadFile(name)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			tb.Fatalf("reading file: %v", err)
		}

		tb.Fatalf("expected file %q to exist, but it does not", name)
	}

	return data
}

func WriteFile(tb testing.TB, data []byte, pathElems ...string) {
	tb.Helper()

	name := filepath.Join(pathElems...)

	if err := os.MkdirAll(filepath.Dir(name), 0o744); err != nil {
		tb.Fatalf("creating file directory: %v", err)
	}

	if err := os.WriteFile(name, data, 0o600); err != nil {
		tb.Fatalf("writing file: %v", err)
	}
}
