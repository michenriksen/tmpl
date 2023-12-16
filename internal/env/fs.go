package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Getwd returns the current working directory.
//
// Works like [os.Getwd] except that it will check the environment variable for
// stubbing the working directory used by tests, and return the value if it's
// defined. This is intended for testing purposes, and should not be used as a
// feature.
func Getwd() (string, error) {
	if dir, ok := LookupEnv(KeyPwd); ok && filepath.IsAbs(dir) {
		return dir, nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current working directory: %w", err)
	}

	return dir, nil
}

// AbsPath returns the absolute path for the given path.
//
// Works like [filepath.Abs] except that it will expand the home directory if
// the path starts with "~".
//
// If the path is already absolute, it will be returned as-is.
func AbsPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("getting user home directory: %w", err)
		}

		return filepath.Join(home, filepath.Clean(path[1:])), nil
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("getting absolute path: %w", err)
	}

	return abs, nil
}
