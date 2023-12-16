package config_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/michenriksen/tmpl/config"
	"github.com/michenriksen/tmpl/internal/testutils"
	"github.com/michenriksen/tmpl/tmux"

	"github.com/stretchr/testify/require"
)

// stubCmd represents an entry for a stubbed command defined in
// testdata/apply-stubcmds.json.
//
// See [loadStubCmds] for more information.
type stubCmd struct {
	Err    error  `json:"err"`
	Output string `json:"output"`
	seen   bool
}

func TestApply(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "project", "cmd"), 0755))

	// Stub HOME and current working directory for consistent test results.
	t.Setenv("HOME", dir)
	t.Setenv("TMPL_PWD", dir)

	cfg, err := config.FromFile(filepath.Join("testdata", "apply.yaml"))
	require.NoError(t, err)

	expectedCmds := loadStubCmds(t)

	var mockCmdRunner tmux.OSCommandRunner = func(_ context.Context, name string, args ...string) ([]byte, error) {
		argStr := strings.Join(args, " ")

		cmd, ok := expectedCmds[argStr]

		if !ok {
			t.Fatalf("unexpected command: %s %s", name, argStr)
		}

		if cmd.seen {
			t.Fatalf("received duplicate command: %s %s", name, argStr)
		}

		t.Logf("received expected command: %s %s", name, argStr)

		cmd.seen = true

		return []byte(cmd.Output), cmd.Err
	}

	cmd, err := tmux.NewRunner(tmux.WithOSCommandRunner(mockCmdRunner))
	require.NoError(t, err)

	session, err := config.Apply(context.Background(), cfg, cmd)
	require.NoError(t, err)

	for args, cmd := range expectedCmds {
		if !cmd.seen {
			t.Fatalf("expected command was not run: %s", args)
		}
	}

	require.Equal(t, "tmpl_test_session", session.Name())
}

// loadStubCmds loads the stub commands defined in testdata/apply-stubcmds.json.
//
// Commands are defined as a map of expected tmux command line arguments mapped
// to a stubCmd struct containing optional stub output and error.
//
// The map keys and stub outputs are expanded using os.ExpandEnv before being
// returned. This allows for using environment variables in the stub commands to
// ensure consistent test results.
func loadStubCmds(t *testing.T) map[string]*stubCmd {
	data := testutils.ReadFile(t, "testdata", "apply-stubcmds.json")

	var cmds map[string]*stubCmd

	if err := json.Unmarshal(data, &cmds); err != nil {
		t.Fatalf("error decoding apply-stubcmds.json: %v", err)
	}

	expanded := make(map[string]*stubCmd, len(cmds))
	for args, cmd := range cmds {
		expanded[os.ExpandEnv(args)] = cmd
	}

	t.Logf("loaded %d stub commands", len(cmds))

	return expanded
}
