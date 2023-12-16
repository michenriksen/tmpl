package tmux_test

import (
	"bytes"
	"context"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/michenriksen/tmpl/internal/testutils"
	"github.com/michenriksen/tmpl/tmux"
)

func TestDefaultRunner_ImplementsRunnerIface(t *testing.T) {
	_, ok := any(&tmux.DefaultRunner{}).(tmux.Runner)

	require.True(t, ok, "expected DefaultRunner to implement Runner interface")
}

func TestDefaultRunner_Run(t *testing.T) {
	tt := []struct {
		name       string
		args       []string
		opts       []tmux.RunnerOption
		assertRun  func(t *testing.T, name string, args ...string) ([]byte, error)
		assertErr  testutils.ErrorAssertion
		wantOutput *regexp.Regexp
	}{
		{
			"command runner success",
			[]string{"new-session", "-d", "-s", "test"},
			[]tmux.RunnerOption{tmux.WithTmux("tmpl_test_tmux")},
			func(t *testing.T, name string, args ...string) ([]byte, error) {
				require.Equal(t, "tmpl_test_tmux", name)
				require.Equal(t, []string{"new-session", "-d", "-s", "test"}, args)

				return []byte{}, nil
			},
			nil, nil,
		},
		{
			"command runner error",
			[]string{},
			nil,
			func(*testing.T, string, ...string) ([]byte, error) {
				return []byte("test output"), exec.ErrNotFound
			},
			testutils.AssertErrorIs(exec.ErrNotFound),
			regexp.MustCompile(`^test output$`),
		},
		{
			"dry-run mode",
			[]string{},
			[]tmux.RunnerOption{tmux.WithDryRunMode(true)},
			func(t *testing.T, _ string, _ ...string) ([]byte, error) {
				t.Fatal("did not expect command runner to be invoked in dry-run mode")
				return []byte{}, nil
			},
			nil, nil,
		},
		{
			"extra options",
			[]string{"test"},
			[]tmux.RunnerOption{tmux.WithTmuxOptions("--extra", "options")},
			func(t *testing.T, _ string, args ...string) ([]byte, error) {
				require.Contains(t, args, "--extra")
				require.Contains(t, args, "options")
				require.Contains(t, args, "test")

				return []byte{}, nil
			}, nil, nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output := bytes.Buffer{}
			logger := slog.New(slog.NewTextHandler(&output, &slog.HandlerOptions{
				Level:       slog.LevelDebug,
				ReplaceAttr: testutils.NewSlogStabilizer(t),
			}))

			cmdRunner := func(_ context.Context, name string, args ...string) ([]byte, error) {
				require.NotNil(t, tc.assertRun, "expected assertRun callback function to be defined")

				t.Log("invoking assertRun callback function")
				return tc.assertRun(t, name, args...)
			}

			execveRunner := func(argv0 string, args, envv []string) error {
				t.Fatalf("unexpected call to Execve: argv0=%q, args=%q, envv=%q", argv0, args, envv)
				return nil
			}

			opts := append([]tmux.RunnerOption{
				tmux.WithLogger(logger),
				tmux.WithOSCommandRunner(cmdRunner),
				tmux.WithSyscallExecRunner(execveRunner),
			}, tc.opts...)

			runner, err := tmux.NewRunner(opts...)
			require.NoError(t, err)

			out, err := runner.Run(context.Background(), tc.args...)

			if tc.assertErr != nil {
				require.Error(t, err)
				tc.assertErr(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, out, "expected non-nil output on success")
			}

			if tc.wantOutput != nil {
				require.Regexp(t, tc.wantOutput, string(out))
			}

			testutils.NewGolden(t).RequireMatch(output.Bytes())
		})
	}
}

func TestDefaultRunner_Run_Integration(t *testing.T) {
	runner, err := tmux.NewRunner(tmux.WithTmux("echo"))
	require.NoError(t, err)

	out, err := runner.Run(context.Background(), "hello world")
	require.NoError(t, err)
	require.Equal(t, "hello world\n", string(out))

	runner, err = tmux.NewRunner(tmux.WithTmux("false"))
	require.NoError(t, err)

	out, err = runner.Run(context.Background())
	require.ErrorContains(t, err, "exit status 1")
	require.Empty(t, out)
}

func TestDefaultRunner_Execve(t *testing.T) {
	stubPath := t.TempDir()
	execPath := filepath.Join(stubPath, "tmux")

	testutils.WriteFile(t, []byte{}, execPath)
	require.NoError(t, os.Chmod(execPath, 0o700))

	t.Setenv("PATH", stubPath)

	tt := []struct {
		name         string
		args         []string
		opts         []tmux.RunnerOption
		assertExecve func(t *testing.T, argv0 string, args, envv []string) error
		assertErr    testutils.ErrorAssertion
	}{
		{
			"syscall exec runner success",
			[]string{"attach-session", "-t", "test"},
			nil,
			func(t *testing.T, argv0 string, args, envv []string) error {
				require.Equal(t, execPath, argv0)
				require.Equal(t, []string{execPath, "attach-session", "-t", "test"}, args)
				require.Equal(t, os.Environ(), envv)

				return nil
			},
			nil,
		},
		{
			"syscall exec runner error",
			[]string{},
			nil,
			func(t *testing.T, argv0 string, args, envv []string) error {
				return fs.ErrNotExist
			},
			testutils.RequireErrorIs(fs.ErrNotExist),
		},
		{
			"dry-run mode",
			[]string{},
			[]tmux.RunnerOption{tmux.WithDryRunMode(true)},
			func(t *testing.T, _ string, _, _ []string) error {
				t.Fatal("did not expect command runner to be invoked in dry-run mode")

				return nil
			},
			nil,
		},
		{
			"extra options",
			[]string{"test"},
			[]tmux.RunnerOption{tmux.WithTmuxOptions("--extra", "options")},
			func(t *testing.T, argv0 string, args, _ []string) error {
				require.Equal(t, argv0, args[0])
				require.Contains(t, args, "--extra")
				require.Contains(t, args, "options")
				require.Contains(t, args, "test")

				return nil
			},
			nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output := bytes.Buffer{}
			logger := slog.New(slog.NewTextHandler(&output, &slog.HandlerOptions{
				Level:       slog.LevelDebug,
				ReplaceAttr: testutils.NewSlogStabilizer(t),
			}))

			execveRunner := func(argv0 string, args, envv []string) error {
				require.NotNil(t, tc.assertExecve, "expected assertExecve callback function to be defined")

				t.Log("invoking assertExecve callback function")
				return tc.assertExecve(t, argv0, args, envv)
			}

			cmdRunner := func(_ context.Context, name string, args ...string) ([]byte, error) {
				t.Fatalf("unexpected call to Run: name=%q, args=%q", name, args)
				return nil, nil
			}

			opts := append([]tmux.RunnerOption{
				tmux.WithLogger(logger),
				tmux.WithOSCommandRunner(cmdRunner),
				tmux.WithSyscallExecRunner(execveRunner),
			}, tc.opts...)

			runner, err := tmux.NewRunner(opts...)
			require.NoError(t, err)

			err = runner.Execve(tc.args...)

			if tc.assertErr != nil {
				require.Error(t, err)
				tc.assertErr(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
