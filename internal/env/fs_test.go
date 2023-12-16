package env_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/michenriksen/tmpl/internal/env"
)

func TestGetwd(t *testing.T) {
	t.Run("no env stub", func(t *testing.T) {
		want, err := os.Getwd()
		require.NoError(t, err)

		got, err := env.Getwd()
		require.NoError(t, err)

		require.Equal(t, want, got)
	})

	t.Run("with env stub", func(t *testing.T) {
		want := t.TempDir()
		t.Setenv("TMPL_PWD", want)

		got, err := env.Getwd()
		require.NoError(t, err)

		require.Equal(t, want, got)
	})
}

func TestAbsPath(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	t.Setenv("HOME", "/home/user")

	tt := []struct {
		name string
		path string
		want string
	}{
		{
			"absolute path",
			"/home/user/project",
			"/home/user/project",
		},
		{
			"relative path",
			"project/scripts",
			filepath.Join(wd, "project/scripts"),
		},
		{
			"relative path traversal",
			"project/scripts/../tests",
			filepath.Join(wd, "project/tests"),
		},
		{
			"relative path dot",
			"./project/scripts",
			filepath.Join(wd, "project/scripts"),
		},
		{
			"tilde path",
			"~/project/scripts",
			"/home/user/project/scripts",
		},
		{
			"tilde path traversal",
			"~/../project/scripts",
			"/home/user/project/scripts",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := env.AbsPath(tc.path)
			require.NoError(t, err)

			require.Equal(t, tc.want, got)
		})
	}
}
