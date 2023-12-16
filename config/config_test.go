package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/michenriksen/tmpl/config"
	"github.com/michenriksen/tmpl/internal/testutils"

	"github.com/stretchr/testify/require"
)

func TestFromFile(t *testing.T) {
	// Stub HOME and current working directory for consistent test results.
	t.Setenv("HOME", "/Users/johndoe")
	t.Setenv("TMPL_PWD", "/Users/johndoe/project")

	tt := []struct {
		name      string
		file      string
		assertErr testutils.ErrorAssertion
	}{
		{"full config", "full.yaml", nil},
		{"minimal config", "minimal.yaml", nil},
		{"tilde home paths", "tilde.yaml", nil},
		{"empty config", "empty.yaml", testutils.RequireErrorIs(config.ErrEmptyConfig)},
		{"broken", "broken.yaml", testutils.RequireErrorContains("decoding error:")},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := config.FromFile(filepath.Join("testdata", tc.file))

			if tc.assertErr != nil {
				require.Error(t, err, "expected error")
				require.Nil(t, cfg, "expected nil config on error")

				tc.assertErr(t, err)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, cfg, "expected non-nil config on success")

			testutils.NewGolden(t).RequireMatch(cfg)
		})
	}
}

func TestFindConfigFile_TraverseDirectories(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "subdir", "second subdir", "third subdir"), 0o744))

	wantCfg := filepath.Join(dir, ".tmpl.yaml")
	testutils.WriteFile(t, []byte("---\nsession:\n  name: test\n"), wantCfg)

	cfg, err := config.FindConfigFile(filepath.Join(dir, "subdir", "second subdir", "third subdir"))

	require.NoError(t, err)
	require.Equal(t, wantCfg, cfg)
}

func TestFindConfigFile_NotFound(t *testing.T) {
	dir := t.TempDir()

	cfg, err := config.FindConfigFile(dir)

	require.ErrorIs(t, err, config.ErrConfigNotFound)
	require.Empty(t, cfg)
}

func TestFindConfigFile_DirNotExist(t *testing.T) {
	cfg, err := config.FindConfigFile("/path/to/non-existent/dir")

	require.ErrorIs(t, err, config.ErrConfigNotFound)
	require.Empty(t, cfg)
}

func TestFindConfigFile_CustomConfigFilename(t *testing.T) {
	t.Setenv("TMPL_CONFIG_NAME", "tmpl-test-custom.yaml")

	dir := t.TempDir()

	wantCfg := filepath.Join(dir, "tmpl-test-custom.yaml")
	testutils.WriteFile(t, []byte("---\nsession:\n  name: test\n"), wantCfg)

	cfg, err := config.FindConfigFile(dir)

	require.NoError(t, err)
	require.Equal(t, wantCfg, cfg)
}
