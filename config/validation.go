package config

import (
	"fmt"
	"regexp"

	"github.com/invopop/validation"

	"github.com/michenriksen/tmpl/internal/rulefuncs"
)

const errorTag = "yaml"

var envVarRE = regexp.MustCompile(`^[A-Z_][A-Z0-9_]+$`)

var nameMatchRule = validation.Match(regexp.MustCompile(`^[\w._-]+$`)).
	Error("must only contain alphanumeric characters, underscores, dots, and dashes")

// Validate validates the configuration.
//
// It checks that:
//
//   - tmux executable exists
//   - session is valid (see [SessionConfig.Validate])
//
// If any of the above checks fail, an error is returned.
func (c Config) Validate() error {
	validation.ErrorTag = errorTag

	return validation.ValidateStruct(&c,
		validation.Field(&c.Tmux, validation.By(rulefuncs.ExecutableExists)),
		validation.Field(&c.Session, validation.Required),
	)
}

// Validate validates the session configuration.
//
// It checks that:
//
//   - session name only contains alphanumeric characters, underscores, dots,
//     and dashes
//   - session path exists
//   - session environment variable names are valid
//   - windows are valid (see [WindowConfig.Validate])
//
// If any of the above checks fail, an error is returned.
func (s SessionConfig) Validate() error {
	validation.ErrorTag = errorTag

	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, nameMatchRule),
		validation.Field(&s.Path, validation.By(rulefuncs.DirExists)),
		validation.Field(&s.Env, validation.By(envVarMapRule)),
		validation.Field(&s.Windows),
	)
}

// Validate validates the window configuration.
//
// It checks that:
//
//   - window name only contains alphanumeric characters, underscores, dots,
//     and dashes
//   - window path exists
//   - window environment variable names are valid
//   - panes are valid (see [PaneConfig.Validate])
//
// If any of the above checks fail, an error is returned.
func (w WindowConfig) Validate() error {
	validation.ErrorTag = errorTag

	return validation.ValidateStruct(&w,
		validation.Field(&w.Name, nameMatchRule),
		validation.Field(&w.Path, validation.By(rulefuncs.DirExists)),
		validation.Field(&w.Env, validation.By(envVarMapRule)),
		validation.Field(&w.Command, validation.Length(1, 0)),
		validation.Field(&w.Commands,
			validation.Each(validation.Length(1, 0)),
		),
		validation.Field(&w.Panes),
	)
}

// Validate validates the pane configuration.
//
// It checks that:
//
//   - pane path exists
//   - pane environment variable names are valid
//   - panes are valid
//
// If any of the above checks fail, an error is returned.
func (p PaneConfig) Validate() error {
	validation.ErrorTag = errorTag

	return validation.ValidateStruct(&p,
		validation.Field(&p.Path, validation.By(rulefuncs.DirExists)),
		validation.Field(&p.Env, validation.By(envVarMapRule)),
		validation.Field(&p.Command, validation.Length(1, 0)),
		validation.Field(&p.Commands,
			validation.Each(validation.Length(1, 0)),
		),
		validation.Field(&p.Panes),
	)
}

// envVarMapRule validates that all keys in a map are valid environment
// variable names (i.e. uppercase letters, numbers and underscores).
func envVarMapRule(val any) error {
	if val == nil {
		return nil
	}

	m, ok := val.(map[string]string)
	if !ok {
		return validation.ErrNotMap
	}

	for k := range m {
		if !envVarRE.MatchString(k) {
			return fmt.Errorf("%q is not a valid environment variable name", k)
		}
	}

	return nil
}
