package cli_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/michenriksen/tmpl/config"
	"github.com/michenriksen/tmpl/internal/cli"
	"github.com/michenriksen/tmpl/internal/testutils"
)

func TestApp_Run_Init(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	stubHome := t.TempDir()
	stubwd := filepath.Join(stubHome, "test project (1)")
	require.NoError(t, os.MkdirAll(stubwd, 0o744))

	t.Setenv("HOME", stubHome)
	t.Setenv("TMPL_PWD", stubwd)

	tt := []struct {
		name        string
		args        []string
		wantCfgPath string
		wantCfg     *config.Config
	}{
		{
			"init current directory",
			[]string{"init"},
			filepath.Join(stubwd, config.ConfigFileName()),
			&config.Config{
				Session: config.SessionConfig{
					Name: "test_project_1",
					Windows: []config.WindowConfig{
						{Name: "main"},
					},
				},
			},
		},
		{
			"init specific directory",
			[]string{"init", filepath.Join(stubHome, "test.project")},
			filepath.Join(stubHome, "test.project", config.ConfigFileName()),
			&config.Config{
				Session: config.SessionConfig{
					Name: "test.project",
					Windows: []config.WindowConfig{
						{Name: "main"},
					},
				},
			},
		},
		{
			"init specific file",
			[]string{"init", filepath.Join(stubHome, "test-project", ".tmpl_config.yml")},
			filepath.Join(stubHome, "test-project", ".tmpl_config.yml"),
			&config.Config{
				Session: config.SessionConfig{
					Name: "test-project",
					Windows: []config.WindowConfig{
						{Name: "main"},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			out := new(bytes.Buffer)

			if tc.wantCfgPath != "" {
				require.NoError(t, os.MkdirAll(filepath.Dir(tc.wantCfgPath), 0o744))
			}

			app, err := cli.NewApp(
				cli.WithOutputWriter(out),
				cli.WithSlogAttrReplacer(testutils.NewSlogStabilizer(t)),
			)
			require.NoError(t, err)

			require.NoError(t, app.Run(ctx, tc.args...))

			testutils.NewGolden(t).RequireMatch(out.Bytes())

			cfg := unmarshalCfg(t, testutils.ReadFile(t, tc.wantCfgPath))

			if tc.wantCfg != nil {
				require.Equal(t, *tc.wantCfg, cfg)
			}
		})
	}
}

func TestApp_Run_InitPlain(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	stubHome := filepath.Join(t.TempDir(), "Zer0_c00l")
	require.NoError(t, os.MkdirAll(stubHome, 0o744))

	t.Setenv("HOME", stubHome)
	t.Setenv("TMPL_PWD", stubHome)

	app, err := cli.NewApp()
	require.NoError(t, err)

	require.NoError(t, app.Run(context.Background(), "init", "-p"))

	wantCfgPath := filepath.Join(stubHome, config.ConfigFileName())
	want := config.Config{
		Session: config.SessionConfig{
			Name: "Zer0_c00l",
			Windows: []config.WindowConfig{
				{Name: "main"},
			},
		},
	}
	got := unmarshalCfg(t, testutils.ReadFile(t, wantCfgPath))

	require.Equal(t, want, got)
	require.NotContains(t, string(testutils.ReadFile(t, wantCfgPath)), "# ",
		"expected configuration to contain no comments",
	)
}

func TestApp_Run_InitFileExists(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	stubHome := t.TempDir()

	t.Setenv("HOME", stubHome)
	t.Setenv("TMPL_PWD", stubHome)

	cfgPath := filepath.Join(stubHome, config.ConfigFileName())
	testutils.WriteFile(t, []byte("don't overwrite me"), cfgPath)

	app, err := cli.NewApp()
	require.NoError(t, err)

	require.NoError(t, app.Run(context.Background(), "init"))

	require.Equal(t, "don't overwrite me", string(testutils.ReadFile(t, cfgPath)),
		"expected configuration file to not be overwritten",
	)
}

func unmarshalCfg(t *testing.T, data []byte) config.Config {
	t.Helper()

	var cfg config.Config

	require.NoError(t, yaml.Unmarshal(data, &cfg), "expected configuration to be parsable YAML")
	require.NoError(t, cfg.Validate(), "expected configuration to be valid")

	return cfg
}
