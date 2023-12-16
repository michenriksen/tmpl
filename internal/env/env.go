// Package env provides functionality for getting information and data from the
// environment.
package env

import "os"

const keyPrefix = "TMPL_"

const (
	// KeyConfigName is the environment variable key for specifying a different
	// configuration file name instead of the default.
	KeyConfigName = "CONFIG_NAME"
	// KeyPwd is the environment variable key for a stubbed working directory
	// used by tests.
	KeyPwd = "PWD"
)

// Getenv retrieves the value of the environment variable named by the key.
//
// Works like [os.Getenv] except that it will prefix the key with an application
// specific prefix to avoid conflicts with other environment variables.
func Getenv(key string) string {
	return os.Getenv(makeKey(key))
}

// LookupEnv retrieves the value of the environment variable named by the key.
//
// Works like [os.LookupEnv] except that it will prefix the key with an
// application specific prefix to avoid conflicts with other environment.
func LookupEnv(key string) (string, bool) {
	return os.LookupEnv(makeKey(key))
}

func makeKey(key string) string {
	return keyPrefix + key
}
