package cli

import (
	"runtime"
	"time"
)

// AppName is the name of the CLI application.
const AppName = "tmpl"

// Build information set by the compiler.
var (
	buildVersion   = "0.0.0-dev"
	buildCommit    = "HEAD"
	buildTime      = ""
	buildGoVersion = runtime.Version()
)

// Version of tmpl.
//
// Returns `0.0.0-dev` if no version is set.
func Version() string {
	return buildVersion
}

// BuildCommit returns the git commit hash tmpl was built from.
//
// Returns `HEAD` if no build commit is set.
func BuildCommit() string {
	return buildCommit
}

// BuildTime returns the UTC time tmpl was built.
//
// Returns current time in UTC if not set.
func BuildTime() string {
	if buildTime == "" {
		return time.Now().UTC().Format(time.RFC3339)
	}

	return buildTime
}

// BuildGoVersion returns the go version tmpl was built with.
//
// Returns version from [runtime.Version] if not set.
func BuildGoVersion() string {
	return buildGoVersion
}
