// Package static provides static embedded assets for the application.
package static

import _ "embed"

//go:embed config.yaml.tmpl
var ConfigTemplate string
