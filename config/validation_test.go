package config_test

import (
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/michenriksen/tmpl/config"
	"github.com/michenriksen/tmpl/internal/testutils"

	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	tt := []struct {
		name      string
		file      string
		assertErr testutils.ErrorAssertion
	}{
		{
			"non-existent tmux executable",
			"invalid-tmux-not-exist.yaml",
			testutils.RequireErrorContains("executable file was not found"),
		},
		{
			"session with invalid name",
			"invalid-session-bad-name.yaml",
			testutils.RequireErrorContains("must only contain alphanumeric characters, underscores, dots, and dashes"),
		},
		{
			"session with non-existent path",
			"invalid-session-path-not-exist.yaml",
			testutils.RequireErrorContains("directory does not exist"),
		},
		{
			"session with invalid env",
			"invalid-session-bad-env.yaml",
			testutils.RequireErrorContains("is not a valid environment variable name"),
		},
		{
			"window with invalid name",
			"invalid-window-bad-name.yaml",
			testutils.RequireErrorContains("must only contain alphanumeric characters, underscores, dots, and dashes"),
		},
		{
			"window with non-existent path",
			"invalid-window-path-not-exist.yaml",
			testutils.RequireErrorContains("directory does not exist"),
		},
		{
			"window with invalid env",
			"invalid-window-bad-env.yaml",
			testutils.RequireErrorContains("is not a valid environment variable name"),
		},
		{
			"pane with non-existent path",
			"invalid-pane-path-not-exist.yaml",
			testutils.RequireErrorContains("directory does not exist"),
		},
		{
			"pane with invalid env",
			"invalid-pane-bad-env.yaml",
			testutils.RequireErrorContains("is not a valid environment variable name"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cfg := loadConfig(t, tc.file)
			err := cfg.Validate()

			if tc.assertErr != nil {
				tc.assertErr(t, err)
				return
			}

			require.NoError(t, err, "expected valid configuration")
		})
	}
}

func loadConfig(t *testing.T, file string) *config.Config {
	t.Helper()

	data := testutils.ReadFile(t, "testdata", file)

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("error decoding content of %s: %v", file, err)
	}

	return &cfg
}
