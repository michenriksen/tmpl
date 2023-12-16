package cli_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/michenriksen/tmpl/internal/cli"
	"github.com/michenriksen/tmpl/internal/mock"
	"github.com/michenriksen/tmpl/internal/testutils"
	"github.com/michenriksen/tmpl/tmux"
)

func TestApp_Run_Apply(t *testing.T) {
	stubHome := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(stubHome, "project", "scripts"), 0o744))

	t.Setenv("NO_COLOR", "1")

	// Stub HOME and current working directory for consistent test results.
	t.Setenv("HOME", stubHome)
	t.Setenv("TMPL_PWD", stubHome)

	// Stub the following environment variables used to determine if the app is
	// running in a tmux session for consistent test results.
	t.Setenv("TMUX", "")
	t.Setenv("TERM_PROGRAM", "")
	t.Setenv("TERM", "xterm-256color")

	dataDir, err := filepath.Abs("testdata")
	require.NoError(t, err)

	// Always include the --debug flag in tests to ensure that the output is
	// included in the golden files.
	alwaysArgs := []string{"--debug"}

	runner, err := tmux.NewRunner()
	require.NoError(t, err)

	stubs := loadTmuxStubs(t)

	tt := []struct {
		name       string
		args       []string
		setupMocks func(*testing.T, *mock.TmuxRunner)
		assertErr  testutils.ErrorAssertion
	}{
		{
			"create new session",
			[]string{"-c", filepath.Join(dataDir, "tmpl.yaml")},
			func(_ *testing.T, r *mock.TmuxRunner) {
				// App gets the current sessions to check if the session already exists.
				stub := stubs["ListSessions"]
				listSess := r.On("Run", stub.Args).Return(stub.Output(), nil).Once()

				// App creates a new session, as it does not exist.
				stub = stubs["NewSession"]
				newSess := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(listSess)

				// App creates the first window named "code".
				stub = stubs["NewWindowCode"]
				newWinCode := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(newSess)

				// App runs on_any hook command in the code window.
				stub = stubs["SendKeysCodeOnAny"]
				codeOnAny := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(newWinCode)

				// App runs on_window hook command in the code window.
				stub = stubs["SendKeysCodeOnWindow"]
				codeOnWindow := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(codeOnAny)

				// App starts Neovim in the "code" window.
				stub = stubs["SendKeysCode"]
				sendKeysCode := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(codeOnWindow)

				// App creates a horizontal pane in the "code" window.
				stub = stubs["NewPaneCode"]
				newPaneCode := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(newWinCode)

				// App runs on_any hook command in the code pane.
				stub = stubs["SendKeysCodePaneOnAny"]
				codePaneOnAny := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(newPaneCode)

				// App runs on_pane hook command in the code pane.
				stub = stubs["SendKeysCodePaneOnPane"]
				codePaneOnPane := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(codePaneOnAny)

				// App starts automatic test run script in the code pane.
				stub = stubs["SendKeysCodePane"]
				r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(codePaneOnPane)

				// App creates the second window named "shell".
				stub = stubs["NewWindowShell"]
				newWinShell := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(newWinCode)

				// App runs on_any hook command in the shell window.
				stub = stubs["SendKeysShellOnAny"]
				shellOnAny := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(newWinShell)

				// App runs on_window hook command in the shell window.
				stub = stubs["SendKeysShellOnWindow"]
				shellOnWindow := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(shellOnAny)

				// App runs `git status` in the shell window.
				stub = stubs["SendKeysShell"]
				r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(shellOnWindow)

				// App creates the third window named "server".
				stub = stubs["NewWindowServer"]
				newWinServer := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(newWinShell)

				// App runs on_any hook command in the server window.
				stub = stubs["SendKeysServerOnAny"]
				serverOnAny := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(newWinServer)

				// App runs on_window hook command in the server window.
				stub = stubs["SendKeysServerOnWindow"]
				serverOnWindow := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(serverOnAny)

				// App starts development server script in the server window.
				stub = stubs["SendKeysServer"]
				r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(serverOnWindow)

				// App creates the fourth window named "prod_logs".
				stub = stubs["NewWindowProdLogs"]
				newWinProdLogs := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(newWinServer)

				// App runs on_any hook command in the prod_logs window.
				stub = stubs["SendKeysProdLogsOnAny"]
				prodLogsOnAny := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(newWinProdLogs)

				// App runs on_window hook command in the prod_logs window.
				stub = stubs["SendKeysProdLogsOnWindow"]
				prodLogsOnWindow := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(prodLogsOnAny)

				// App connects to production host via SSH in the prod_logs window.
				stub = stubs["SendKeysProdLogsSSH"]
				prodLogsSSH := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(prodLogsOnWindow)

				// App navigates to the logs directory in the prod_logs window.
				stub = stubs["SendKeysProdLogsCdLogs"]
				prodLogsCdLogs := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(prodLogsSSH)

				// App tails the application log file in the prod_logs window.
				stub = stubs["SendKeysProdLogsTail"]
				r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(prodLogsCdLogs)

				// App selects the code window.
				stub = stubs["SelectWindowCode"]
				selectWinCode := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(sendKeysCode)

				// App determines the pane base index to select the initial pane.
				stub = stubs["PaneBaseIndexOpt"]
				paneBaseIndexOpt := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(selectWinCode)

				// App selects the initial code pane running Neovim.
				stub = stubs["SelectPaneCode"]
				selectPaneCode := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(paneBaseIndexOpt)

				// Finally, App attaches the new session.
				stub = stubs["AttachSession"]
				r.On("Execve", stub.Args).Return(nil).Once().NotBefore(selectPaneCode)
			},
			nil,
		},
		{
			"session exists",
			[]string{"-c", filepath.Join(dataDir, "tmpl.yaml")},
			func(_ *testing.T, r *mock.TmuxRunner) {
				// App gets the current sessions to check if the session already exists.
				stub := stubs["ListSessionsExists"]
				r.On("Run", stub.Args).Return(stub.Output(), nil).Once()

				// Since the session already exists, it attaches it.
				stub = stubs["AttachSession"]
				r.On("Execve", stub.Args).Return(nil).Once()
			},
			nil,
		},
		{
			"new session fails",
			[]string{"-c", filepath.Join(dataDir, "tmpl.yaml")},
			func(_ *testing.T, r *mock.TmuxRunner) {
				// App gets the current sessions to check if the session already exists.
				stub := stubs["ListSessions"]
				listSess := r.On("Run", stub.Args).Return(stub.Output(), nil).Once()

				// App creates a new session but it fails.
				stub = stubs["NewSession"]
				r.On("Run", stub.Args).
					Return([]byte("failed to connect to server: Connection refused"), errors.New("exit status 1")).Once().NotBefore(listSess)
			},
			testutils.RequireErrorContains("running new-session command: exit status 1"),
		},
		{
			"new window fails",
			[]string{"-c", filepath.Join(dataDir, "tmpl.yaml")},
			func(_ *testing.T, r *mock.TmuxRunner) {
				// App gets the current sessions to check if the session already exists.
				stub := stubs["ListSessions"]
				listSess := r.On("Run", stub.Args).Return(stub.Output(), nil).Once()

				// App creates a new session, as it does not exist.
				stub = stubs["NewSession"]
				newSess := r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(listSess)

				// App creates the first window named "code" but it fails.
				stub = stubs["NewWindowCode"]
				newWinCode := r.On("Run", stub.Args).
					Return([]byte("failed to connect to server: Connection refused"), errors.New("exit status 1")).Once().NotBefore(newSess)

				// App closes the failed session.
				stub = stubs["CloseSession"]
				r.On("Run", stub.Args).Return(stub.Output(), nil).Once().NotBefore(newWinCode)
			},
			testutils.RequireErrorContains("running new-window command: exit status 1"),
		},
		{
			"broken config file",
			[]string{"-c", filepath.Join(dataDir, "tmpl-broken.yaml")},
			nil,
			testutils.RequireErrorContains("invalid configuration"),
		},
		{
			"invalid config file",
			[]string{"-c", filepath.Join(dataDir, "tmpl-invalid.yaml")},
			nil,
			testutils.RequireErrorContains("invalid configuration"),
		},
		{
			"show help",
			[]string{"-h"},
			nil,
			testutils.RequireErrorIs(cli.ErrHelp),
		},
		{
			"show version",
			[]string{"-v"},
			nil,
			testutils.RequireErrorIs(cli.ErrVersion),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			mockRunner := mock.NewTmuxRunner(t, runner)

			if tc.setupMocks != nil {
				tc.setupMocks(t, mockRunner)
			}

			app, err := cli.NewApp(
				cli.WithOutputWriter(out),
				cli.WithTmux(mockRunner),
				cli.WithSlogAttrReplacer(testutils.NewSlogStabilizer(t)),
			)
			require.NoError(t, err)

			args := append(alwaysArgs, tc.args...)

			err = app.Run(context.Background(), args...)

			if tc.assertErr != nil {
				require.Error(t, err)
				tc.assertErr(t, err)
			} else {
				require.NoError(t, err)
			}

			if !mockRunner.AssertExpectations(t) {
				t.FailNow()
			}

			testutils.NewGolden(t).RequireMatch(testutils.Stabilize(t, out.Bytes()))
		})
	}
}

// loadTmuxStubs loads the expected tmux command arguments and stub output from
// the tmux-stubs.yaml file in the testdata directory.
func loadTmuxStubs(t *testing.T) map[string]tmuxStub {
	t.Helper()

	data := testutils.ReadFile(t, "testdata", "tmux-stubs.yaml")

	stubs := make(map[string]tmuxStub)

	require.NoError(t, yaml.Unmarshal(data, stubs), "expected tmux-stubs.yaml to contain valid YAML")

	return stubs
}

// tmuxStub contains expected tmux command arguments and stub output to use in
// tests.
type tmuxStub struct {
	OutputString string   `yaml:"output"`
	Args         []string `yaml:"args"`
}

// Output returns the stub output as a byte slice.
func (s tmuxStub) Output() []byte {
	return []byte(s.OutputString)
}
