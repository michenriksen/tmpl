package rulefuncs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/michenriksen/tmpl/internal/rulefuncs"
	"github.com/michenriksen/tmpl/internal/testutils"
)

func TestDirExists(t *testing.T) {
	dir := t.TempDir()

	file, err := os.CreateTemp(dir, "test")
	require.NoError(t, err)

	tt := []struct {
		name      string
		val       any
		assertErr testutils.ErrorAssertion
	}{
		{
			"empty string",
			"",
			nil,
		},
		{
			"nil value",
			nil,
			nil,
		},
		{
			"existing directory",
			dir,
			nil,
		},
		{
			"non-existing directory",
			"/tmpl/test/not-exist",
			testutils.AssertErrorContains("directory does not exist"),
		},
		{
			"file path",
			file.Name(),
			testutils.AssertErrorContains("not a directory"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := rulefuncs.DirExists(tc.val)

			if tc.assertErr != nil {
				tc.assertErr(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestFileExists(t *testing.T) {
	dir := t.TempDir()

	file, err := os.CreateTemp(dir, "test")
	require.NoError(t, err)

	tt := []struct {
		name      string
		val       any
		assertErr testutils.ErrorAssertion
	}{
		{
			"empty string",
			"",
			nil,
		},
		{
			"nil value",
			nil,
			nil,
		},
		{
			"existing file",
			file.Name(),
			nil,
		},
		{
			"non-existing file",
			"/tmpl/test/not-exist",
			testutils.AssertErrorContains("file does not exist"),
		},
		{
			"directory path",
			dir,
			testutils.AssertErrorContains("not a file"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := rulefuncs.FileExists(tc.val)

			if tc.assertErr != nil {
				tc.assertErr(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestExecutableExists(t *testing.T) {
	dir := t.TempDir()

	execPath := filepath.Join(dir, "tmpl_test_exec")
	nonExecPath := filepath.Join(dir, "tmpl_test_non_exec")

	testutils.WriteFile(t, []byte{}, execPath)
	testutils.WriteFile(t, []byte{}, nonExecPath)
	os.Chmod(execPath, 0o700)

	t.Setenv("PATH", dir)

	tt := []struct {
		name      string
		val       any
		assertErr testutils.ErrorAssertion
	}{
		{
			"empty string",
			"",
			nil,
		},
		{
			"nil value",
			nil,
			nil,
		},
		{
			"executable absolute path",
			execPath,
			nil,
		},
		{
			"executable name",
			"tmpl_test_exec",
			nil,
		},
		{
			"non-executable absolute path",
			nonExecPath,
			testutils.AssertErrorContains("permission denied"),
		},
		{
			"non-executable name",
			"tmpl_test_non_exec",
			testutils.AssertErrorContains("executable file not found in $PATH"),
		},
		{
			"non-existing absolute path",
			"/tmpl/test/not-exist",
			testutils.AssertErrorContains("executable file was not found"),
		},
		{
			"directory path",
			dir,
			testutils.AssertErrorContains("is a directory"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := rulefuncs.ExecutableExists(tc.val)

			if tc.assertErr != nil {
				tc.assertErr(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
