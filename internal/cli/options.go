package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"text/template"
)

const globalOpts = `Global options:

    -d, --debug                enable debug logging
    -h, --help                 show this message and exit
    -j, --json                 enable JSON logging
    -q, --quiet                enable quiet logging
    -v, --version              show the version and exit`

const subCmds = `Available commands:

    apply (default)            apply configuration and attach session
    check                      validate configuration file
    init                       generate a new configuration file`

const usageTmpl = `Usage: {{ .AppName }} [command] [options] [args]

Simple tmux session management.

{{ .Commands }}

{{ .GlobalOptions }}

Examples:

    # apply nearest configuration file and attach/switch client to session:
    $ {{ .AppName }}

    # or explicitly:
    $ {{ .AppName }} -c /path/to/config.yaml

    # generate a new configuration file in the current working directory:
    $ {{ .AppName }} init
`

const applyUsageTmpl = `Usage: {{ .AppName }} apply [options]

Creates a new tmux session from a {{ .AppName }} configuration file and then
connects to it.

If the session already exists, the configuration process is skipped.


Options:

    -c, --config PATH          configuration file path (default: find nearest)
    -n, --dry-run              enable dry-run mode

{{ .GlobalOptions }}

Examples:

    # apply nearest configuration file and attach/switch client to session:
    $ {{ .AppName }} apply

    # or explicitly:
    $ {{ .AppName }} apply -c /path/to/config.yaml

    # simulate applying configuration file. No tmux commands are executed:
    $ {{ .AppName }} apply --dry-run
`

const initUsageTmpl = `Usage: {{ .AppName }} init [options] [path]

Generates a skeleton {{ .AppName }} configuration file to get you started.


Options:
    -p, --plain                make plain configuration with no comments

{{ .GlobalOptions }}

Examples:

    # create a configuration file in the current working directory:
    $ {{ .AppName }} init

    # or at a specific location:
    $ {{ .AppName }} init /path/to/config.yaml
`

const checkUsageTmpl = `Usage: {{ .AppName }} check [options] [path]

Performs validation of a {{ .AppName }} configuration file and reports whether
it is valid or not.


Options:

    -c, --config PATH          configuration file path (default: find nearest)

{{ .GlobalOptions }}

Examples:

    # validate configuration file in the current working directory:
    $ {{ .AppName }} check

    # or at a specific location:
    $ {{ .AppName }} check -c /path/to/config.yaml
`

const versionTmpl = `{{ .AppName }}:
  Version:    {{ .Version }}
  Go version: {{ .GoVersion }}
  Git commit: {{ .Commit }}
  Released:   {{ .BuildTime }}
`

var (
	ErrHelp    = errors.New("help requested")
	ErrVersion = errors.New("version requested")
)

// options represents the command-line options for the CLI application.
type options struct {
	args    []string
	version bool
	help    bool

	// Global options.
	Debug bool
	Quiet bool
	JSON  bool

	// Options for apply sub-command.
	ConfigPath string
	DryRun     bool

	// Options for init sub-command.
	Plain bool
}

// parseApplyOptions parses the command-line options for the apply sub-command.
func parseApplyOptions(args []string, output io.Writer) (*options, error) {
	flagSet := flag.NewFlagSet("apply", flag.ContinueOnError)
	isSubCmd := len(args) != 0 && args[0] == "apply"

	flagSet.SetOutput(output)
	flagSet.Usage = func() {
		var (
			usage string
			err   error
		)

		if isSubCmd {
			usage, err = renderOptsTemplate(applyUsageTmpl)
		} else {
			usage, err = renderOptsTemplate(usageTmpl)
		}

		if err != nil {
			panic(err)
		}

		fmt.Fprint(output, usage)
	}

	opts := &options{}
	initGlobalOpts(flagSet, opts)

	flagSet.StringVar(&opts.ConfigPath, "config", "", "path to the configuration file")
	flagSet.StringVar(&opts.ConfigPath, "c", "", "path to the configuration file")
	flagSet.BoolVar(&opts.DryRun, "dry-run", false, "enable dry-run mode")
	flagSet.BoolVar(&opts.DryRun, "n", false, "enable dry-run mode")

	if isSubCmd {
		args = args[1:]
	}

	opts, err := parseFlagSet(args, flagSet, opts)
	if err != nil {
		return nil, err
	}

	if len(opts.args) != 0 {
		return nil, fmt.Errorf("unknown command: %s", opts.args[0])
	}

	return opts, nil
}

// parseInitOptions parses the command-line options for the init sub-command.
func parseInitOptions(args []string, output io.Writer) (*options, error) {
	flagSet := flag.NewFlagSet("init", flag.ContinueOnError)

	flagSet.SetOutput(output)
	flagSet.Usage = func() {
		usage, err := renderOptsTemplate(initUsageTmpl)
		if err != nil {
			panic(err)
		}

		fmt.Fprint(output, usage)
	}

	opts := &options{}
	initGlobalOpts(flagSet, opts)

	flagSet.BoolVar(&opts.Plain, "plain", false, "make plain configuration with no comments")
	flagSet.BoolVar(&opts.Plain, "p", false, "make plain configuration with no comments")

	return parseFlagSet(args, flagSet, opts)
}

func parseCheckOptions(args []string, output io.Writer) (*options, error) {
	flagSet := flag.NewFlagSet("check", flag.ContinueOnError)

	flagSet.SetOutput(output)
	flagSet.Usage = func() {
		usage, err := renderOptsTemplate(checkUsageTmpl)
		if err != nil {
			panic(err)
		}

		fmt.Fprint(output, usage)
	}

	opts := &options{}
	initGlobalOpts(flagSet, opts)

	flagSet.StringVar(&opts.ConfigPath, "config", "", "path to the configuration file")
	flagSet.StringVar(&opts.ConfigPath, "c", "", "path to the configuration file")

	return parseFlagSet(args, flagSet, opts)
}

func initGlobalOpts(flagSet *flag.FlagSet, opts *options) {
	flagSet.BoolVar(&opts.Debug, "debug", false, "enable debug logging")
	flagSet.BoolVar(&opts.Debug, "d", false, "enable debug logging")
	flagSet.BoolVar(&opts.Quiet, "quiet", false, "enable quiet logging")
	flagSet.BoolVar(&opts.Quiet, "q", false, "enable quiet logging")
	flagSet.BoolVar(&opts.JSON, "json", false, "enable JSON logging")
	flagSet.BoolVar(&opts.JSON, "j", false, "enable JSON logging")
	flagSet.BoolVar(&opts.version, "version", false, "show the version and exit")
	flagSet.BoolVar(&opts.version, "v", false, "show the version and exit")
	flagSet.BoolVar(&opts.help, "help", false, "show this message and exit")
	flagSet.BoolVar(&opts.help, "h", false, "show this message and exit")
}

func parseFlagSet(args []string, flagSet *flag.FlagSet, opts *options) (*options, error) {
	if err := flagSet.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}

	if opts.help {
		flagSet.Usage()
		return nil, ErrHelp
	}

	if opts.version {
		info, err := renderOptsTemplate(versionTmpl)
		if err != nil {
			return nil, err
		}

		fmt.Fprint(flagSet.Output(), info)

		return nil, ErrVersion
	}

	opts.args = flagSet.Args()

	return opts, nil
}

func renderOptsTemplate(s string) (string, error) {
	data := optsTemplateData{
		AppName:       AppName,
		BuildTime:     BuildTime(),
		Commands:      subCmds,
		Commit:        BuildCommit(),
		GlobalOptions: globalOpts,
		GoVersion:     BuildGoVersion(),
		Version:       Version(),
	}

	tmpl, err := template.New("usage").Parse(s)
	if err != nil {
		return "", fmt.Errorf("parsing options template: %w", err)
	}

	b := strings.Builder{}

	if err := tmpl.Execute(&b, data); err != nil {
		return "", fmt.Errorf("rendering options template: %w", err)
	}

	return b.String(), nil
}

type optsTemplateData struct {
	AppName       string
	BuildTime     string
	Commands      string
	Commit        string
	GlobalOptions string
	GoVersion     string
	Version       string
}
