// Package rulefuncs provides custom validation rule functions.
package rulefuncs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/invopop/validation"
)

// DirExists validated that provided value is a valid path to an existing
// directory.
func DirExists(value any) error {
	if value == nil {
		return nil
	}

	dir, err := validation.EnsureString(value)
	if err != nil {
		return err
	}

	if dir == "" {
		return nil
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		return errors.New("invalid directory path")
	}

	info, err := os.Stat(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return errors.New("directory does not exist")
		}

		return validation.NewInternalError(fmt.Errorf("getting directory info: %w", err))
	}

	if !info.IsDir() {
		return errors.New("not a directory")
	}

	return nil
}

// FileExists validates that provided value is a valid path to an existing file.
func FileExists(value any) error {
	if value == nil {
		return nil
	}

	name, err := validation.EnsureString(value)
	if err != nil {
		return err
	}

	if name == "" {
		return nil
	}

	name, err = filepath.Abs(name)
	if err != nil {
		return errors.New("invalid file path")
	}

	info, err := os.Stat(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return errors.New("file does not exist")
		}

		return validation.NewInternalError(fmt.Errorf("getting file info: %w", err))
	}

	if info.IsDir() {
		return errors.New("not a file")
	}

	return nil
}

// ExecutableExists validates that provided value is an executable command
// in PATH or an absolute path to an executable file.
func ExecutableExists(value any) error {
	if value == nil {
		return nil
	}

	name, err := validation.EnsureString(value)
	if err != nil {
		return err
	}

	if name == "" {
		return nil
	}

	if _, err := exec.LookPath(name); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return errors.New("executable file was not found")
		}

		return validation.NewInternalError(fmt.Errorf("looking up command: %w", err))
	}

	return nil
}
