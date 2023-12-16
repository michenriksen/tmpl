// Package mock provides mock implementations of components to use in tests.
package mock

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/michenriksen/tmpl/tmux"
)

type TmuxRunner struct {
	tb      testing.TB
	wrapped tmux.Runner

	mock.Mock
}

func NewTmuxRunner(tb testing.TB, wrapped tmux.Runner) *TmuxRunner {
	tb.Helper()

	return &TmuxRunner{tb: tb, wrapped: wrapped}
}

func (r *TmuxRunner) Run(ctx context.Context, args ...string) ([]byte, error) {
	r.tb.Helper()

	ret := r.Called(cleanArgs(args))

	if ret.Is(nil) {
		r.tb.Log("Run call mock has nil return; delegating to wrapped runner")
		return r.wrapped.Run(ctx, args...) //nolint:wrapcheck // Intentional.
	}

	output, ok := ret.Get(0).([]byte)
	if !ok {
		r.tb.Fatalf("Run call mock has invalid output return type: %T", ret.Get(0))
	}

	return output, ret.Error(1) //nolint:wrapcheck // Intentional.
}

func (r *TmuxRunner) Execve(args ...string) error {
	r.tb.Helper()

	return r.Called(cleanArgs(args)).Error(0) //nolint:wrapcheck // Intentional.
}

func (r *TmuxRunner) IsDryRun() bool {
	r.tb.Helper()
	return false
}

func (r *TmuxRunner) Debug(msg string, args ...any) {
	r.tb.Helper()
	r.wrapped.Debug(msg, append(args, "mock", true)...)
}

func (r *TmuxRunner) Log(msg string, args ...any) {
	r.tb.Helper()
	r.wrapped.Log(msg, append(args, "mock", true)...)
}

func (r *TmuxRunner) SetLogger(logger *slog.Logger) {
	r.tb.Helper()
	r.wrapped.SetLogger(logger)
}

func cleanArgs(args []string) []string {
	tmpDir := os.TempDir()
	res := make([]string, 0, len(args))

	for _, arg := range args {
		if !strings.HasPrefix(arg, tmpDir) {
			res = append(res, arg)
			continue
		}

		res = append(res, "/tmp/path")
	}

	return res
}
