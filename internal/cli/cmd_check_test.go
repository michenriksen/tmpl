package cli_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/michenriksen/tmpl/config"
	"github.com/michenriksen/tmpl/internal/cli"
	"github.com/michenriksen/tmpl/internal/testutils"
)

func TestApp_Run_Check(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	stubHome := t.TempDir()
	t.Setenv("HOME", stubHome)
	t.Setenv("TMPL_PWD", stubHome)

	require.NoError(t, os.MkdirAll(filepath.Join(stubHome, "project", "scripts"), 0o744))

	testutils.WriteFile(t,
		testutils.ReadFile(t, "testdata", "tmpl.yaml"),
		stubHome, config.ConfigFileName(),
	)

	testutils.WriteFile(t,
		testutils.ReadFile(t, "testdata", "tmpl.yaml"),
		stubHome, "project", config.ConfigFileName(),
	)

	testutils.WriteFile(t,
		testutils.ReadFile(t, "testdata", "tmpl-broken.yaml"),
		stubHome, ".tmpl.broken.yaml",
	)

	testutils.WriteFile(t,
		testutils.ReadFile(t, "testdata", "tmpl-invalid.yaml"),
		stubHome, ".tmpl.invalid.yaml",
	)

	tt := []struct {
		name      string
		args      []string
		assertErr testutils.ErrorAssertion
	}{
		{
			"check current config",
			[]string{"check"},
			nil,
		},
		{
			"check specific config",
			[]string{"check", "-c", filepath.Join(stubHome, "project", config.ConfigFileName())},
			nil,
		},
		{
			"unparsable config",
			[]string{"check", "-c", filepath.Join(stubHome, ".tmpl.broken.yaml")},
			testutils.RequireErrorIs(cli.ErrInvalidConfig),
		},
		{
			"invalid config",
			[]string{"check", "-c", filepath.Join(stubHome, ".tmpl.invalid.yaml")},
			testutils.RequireErrorIs(cli.ErrInvalidConfig),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			out := new(bytes.Buffer)

			app, err := cli.NewApp(
				cli.WithOutputWriter(out),
				cli.WithSlogAttrReplacer(testutils.NewSlogStabilizer(t)),
			)
			require.NoError(t, err)

			err = app.Run(context.Background(), tc.args...)

			if tc.assertErr != nil {
				require.Error(t, err)
				tc.assertErr(t, err)
			} else {
				require.NoError(t, err)
			}

			testutils.NewGolden(t).RequireMatch(out.Bytes())
		})
	}
}
