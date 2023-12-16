package testutils

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

var (
	tempDir = os.TempDir()
	homeDir = os.Getenv("HOME")
)

var replacers = []struct {
	re   *regexp.Regexp
	repl []byte
}{
	{
		// RFC3339 timestamps.
		regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z\b`),
		[]byte("0001-01-01T00:00:00Z"),
	},
	{
		// Go version strings.
		regexp.MustCompile(`\bgo1\.\d+\.\d+\b`),
		[]byte("go0.0.0"),
	},
	{
		// Temporary directory paths.
		regexp.MustCompile(fmt.Sprintf(`\b\/%s\/[\w\/_-]+\b`, regexp.QuoteMeta(tempDir))),
		[]byte("/tmp/path"),
	},
	{
		// Home directory paths.
		regexp.MustCompile(fmt.Sprintf(`\b\/%s\/[\w\/_-]+\b`, regexp.QuoteMeta(homeDir))),
		[]byte("/home/user"),
	},
}

// Stabilize replaces non-deterministic and environment-dependent values in
// [data] with stable values.
//
// The function replaces all timestamps in RFC3339 format with the zero value
// and all Go version strings with "go0.0.0".
func Stabilize(tb testing.TB, data []byte) []byte {
	tb.Helper()

	for _, r := range replacers {
		data = r.re.ReplaceAll(data, r.repl)
	}

	return data
}

// NewSlogStabilizer returns an attribute replacer function for [slog.Logger]
// that replaces non-deterministic and environment-dependent attribute values
// with stable values.
//
// The function replaces all attribute values of type [time.Time] and
// [time.Duration] with their zero values.
//
// If a string value begins with the system's temporary directory path, or the
// current user's home directory, the path is replaced with a stable path.
//
// Usage:
//
//	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
//	    ReplaceAttr: testutils.NewSlogStabilizer(t)
//	}))
func NewSlogStabilizer(tb testing.TB) func([]string, slog.Attr) slog.Attr {
	tb.Helper()

	return func(_ []string, a slog.Attr) slog.Attr {
		tb.Helper()

		switch a.Value.Any().(type) {
		case time.Time:
			return slog.Attr{Key: a.Key, Value: slog.TimeValue(time.Time{})}
		case time.Duration:
			return slog.Attr{Key: a.Key, Value: slog.DurationValue(time.Duration(0))}
		case string:
			val := a.Value.String()

			if strings.HasPrefix(val, string(filepath.Separator)) && isPathKey(a.Key) {
				val = filepath.Join("/stabilized/path", filepath.Base(val))
				return slog.Attr{Key: a.Key, Value: slog.StringValue(val)}
			}

			return slog.Attr{Key: a.Key, Value: slog.StringValue(val)}
		default:
			return a
		}
	}
}

func isPathKey(key string) bool {
	for _, pk := range []string{"path", "dir", "file"} {
		if strings.Contains(key, pk) {
			return true
		}
	}

	return false
}
