package config

import (
	"errors"
	"fmt"
)

var (
	// ErrConfigNotFound is returned when a configuration file is not found.
	ErrConfigNotFound = errors.New("configuration file not found")
	// ErrEmptyConfig is returned when a configuration file is empty.
	ErrEmptyConfig = errors.New("configuration file is empty")
	// ErrInvalidConfig is returned when a configuration file contains invalid
	// and unparsable YAML.
	ErrInvalidConfig = errors.New("configuration file is not parsable")
)

// DecodeError is returned when a configuration file cannot be decoded.
type DecodeError struct {
	err  error
	path string
}

func decodeError(err error, path string) DecodeError {
	return DecodeError{err: err, path: path}
}

// Error implements the error interface.
func (e DecodeError) Error() string {
	return fmt.Sprintf("decoding error: %s", e.err)
}

// Unwrap implements the [errors.Wrapper] interface.
func (e DecodeError) Unwrap() error {
	return e.err
}

// Path returns the path to the configuration file that was attempted to be
// decoded.
func (e DecodeError) Path() string {
	return e.path
}
