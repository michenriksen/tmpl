package env_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/michenriksen/tmpl/internal/env"
)

func TestGetenv(t *testing.T) {
	t.Setenv("TMPL_TEST_GET_ENV", "good")
	t.Setenv("TEST_GET_ENV", "bad")

	require.Equal(t, "good", env.Getenv("TEST_GET_ENV"))
}

func TestLookupEnv(t *testing.T) {
	t.Setenv("TMPL_TEST_GET_ENV", "good")
	t.Setenv("TEST_GET_ENV", "bad")

	val, ok := env.LookupEnv("TEST_GET_ENV")
	require.True(t, ok)
	require.Equal(t, "good", val)
}
