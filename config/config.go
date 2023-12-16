// Package config loads, validates, and applies tmpl configurations.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"

	"github.com/michenriksen/tmpl/internal/env"
)

// DefaultConfigFile is the default configuration filename.
const DefaultConfigFile = ".tmpl.yaml"

var specialCharsRegexp = regexp.MustCompile(`[^\w_]+`)

// Config represents a session configuration loaded from a YAML file.
type Config struct {
	path        string
	Session     SessionConfig `yaml:"session"`      // Session configuration.
	Tmux        string        `yaml:"tmux"`         // Path to tmux executable.
	TmuxOptions []string      `yaml:"tmux_options"` // Additional tmux options.
}

// FromFile loads a session configuration from provided file path.
//
// File is expected to be in YAML format.
func FromFile(cfgPath string) (*Config, error) {
	cfg, err := load(cfgPath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Path returns the path to the configuration file from which the configuration
// was loaded.
func (c *Config) Path() string {
	return c.path
}

// NumWindows returns the number of window configurations for the session.
func (c *Config) NumWindows() int {
	n := len(c.Session.Windows)
	if n == 0 {
		return 1
	}

	return n
}

// NumPanes returns the number of pane configurations for the session.
func (c *Config) NumPanes() int {
	n := 0

	for _, w := range c.Session.Windows {
		n += len(w.Panes)

		for _, p := range w.Panes {
			n += len(p.Panes)
		}
	}

	return n
}

// SessionConfig represents a tmux session configuration. It contains the name
// of the session, the path to the directory where the session will be created
// and the window configurations.
//
// Any environment variables defined in the session configuration will be
// inherited by all windows and panes.
type SessionConfig struct {
	Name     string            `yaml:"name"`      // Session name.
	Path     string            `yaml:"path"`      // Session directory.
	OnWindow string            `yaml:"on_window"` // Shell command to run in all windows.
	OnPane   string            `yaml:"on_pane"`   // Shell command to run in all panes.
	OnAny    string            `yaml:"on_any"`    // Shell command to run in all windows and panes.
	Env      map[string]string `yaml:"env"`       // Session environment variables.
	Windows  []WindowConfig    `yaml:"windows"`   // Window configurations.
}

// WindowConfig represents a tmux window configuration. It contains the name of
// the window, the path to the directory where the window will be created, the
// command to run in the window and pane configurations.
//
// If a path is not specified, a window will inherit the session path.
//
// Any environment variables defined in the window configuration will be
// inherited by all panes. If a variable is defined in both the session and
// window configuration, the window variable will take precedence.
type WindowConfig struct {
	Name     string            `yaml:"name"`     // Window name.
	Path     string            `yaml:"path"`     // Window directory.
	Command  string            `yaml:"command"`  // Command to run in the window.
	Commands []string          `yaml:"commands"` // Commands to run in the window.
	Env      map[string]string `yaml:"env"`      // Window environment variables.
	Panes    []PaneConfig      `yaml:"panes"`    // Pane configurations.
	Active   bool              `yaml:"active"`   // Whether the window should be selected.
}

// PaneConfig represents a tmux pane configuration. It contains the path to the
// directory where the pane will be created, the command to run in the pane,
// the size of the pane and whether the pane should be split horizontally or
// vertically.
//
// If a path is not specified, a pane will inherit the window path.
//
// Any inherited environment variables from the window or session will be
// overridden by variables defined in the pane configuration if they have the
// same name.
type PaneConfig struct {
	Env        map[string]string `yaml:"env"`        // Pane environment variables.
	Path       string            `yaml:"path"`       // Pane directory.
	Command    string            `yaml:"command"`    // Command to run in the pane.
	Commands   []string          `yaml:"commands"`   // Commands to run in the pane.
	Size       string            `yaml:"size"`       // Pane size (cells or percentage)
	Horizontal bool              `yaml:"horizontal"` // Whether the pane should be split horizontally.
	Panes      []PaneConfig      `yaml:"panes"`      // Pane configurations.
	Active     bool              `yaml:"active"`     // Whether the pane should be selected.
}

// FindConfigFile searches for a configuration file starting from the provided
// directory and going up until the root directory is reached. If no file is
// found, ErrConfigNotFound is returned.
//
// By default, the configuration file name is .tmpl.yaml. This can be changed
// by setting the TMPL_CONFIG_FILE environment variable.
func FindConfigFile(dir string) (string, error) {
	name := ConfigFileName()

	for {
		cfgPath := filepath.Join(dir, name)

		info, err := os.Stat(cfgPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				if dir == "/" {
					return "", ErrConfigNotFound
				}

				dir = filepath.Dir(dir)

				continue
			}

			return "", fmt.Errorf("getting file info: %w", err)
		}

		if info.IsDir() {
			return "", fmt.Errorf("path %q is a directory", cfgPath)
		}

		return cfgPath, nil
	}
}

// load reads and decodes a YAML configuration file into a Config struct and
// sets default values.
func load(cfgPath string) (*Config, error) {
	var cfg Config

	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("opening configuration file: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("getting configuration file info: %w", err)
	}

	if info.Size() == 0 {
		return nil, ErrEmptyConfig
	}

	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)

	if err := decoder.Decode(&cfg); err != nil {
		return nil, decodeError(err, cfgPath)
	}

	cfg.path = cfgPath

	if err := setDefaults(&cfg); err != nil {
		return nil, fmt.Errorf("setting default values: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default values for blank configuration fields.
//
// The following fields are set to default values:
//
// - Session.Name: defaults to <current directory name>.
// - Session.Path: defaults to current working directory.
// - Window.Path: defaults to Session.Path.
// - Pane.Path: defaults to Window.Path.
func setDefaults(cfg *Config) error {
	wd, err := env.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %w", err)
	}

	if cfg.Session.Name == "" {
		name := specialCharsRegexp.ReplaceAllString(filepath.Base(wd), "_")
		cfg.Session.Name = name
	}

	if cfg.Session.Path == "" {
		cfg.Session.Path = wd
	} else {
		if cfg.Session.Path, err = env.AbsPath(cfg.Session.Path); err != nil {
			return fmt.Errorf("expanding session path: %w", err)
		}
	}

	for i, w := range cfg.Session.Windows {
		if w.Path == "" {
			w.Path = cfg.Session.Path
		} else {
			if w.Path, err = env.AbsPath(w.Path); err != nil {
				return fmt.Errorf("expanding window path: %w", err)
			}
		}

		for j, p := range w.Panes {
			if p.Path == "" {
				p.Path = w.Path
			} else {
				if p.Path, err = env.AbsPath(p.Path); err != nil {
					return fmt.Errorf("expanding pane path: %w", err)
				}
			}

			w.Panes[j] = p
		}

		cfg.Session.Windows[i] = w
	}

	return nil
}

// ConfigFileName returns the name of the configuration that.
//
// Returns the value of the TMPL_CONFIG_NAME environment variable if set,
// otherwise it returns [DefaultConfigFile].
func ConfigFileName() string {
	if name := env.Getenv(env.KeyConfigName); name != "" {
		return name
	}

	return DefaultConfigFile
}
