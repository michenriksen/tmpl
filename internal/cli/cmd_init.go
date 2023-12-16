package cli

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/michenriksen/tmpl/config"
	"github.com/michenriksen/tmpl/internal/env"
	"github.com/michenriksen/tmpl/internal/static"
)

var cleanSessNameRE = regexp.MustCompile(`[^\w._-]+`)

func (a *App) runInit(_ context.Context) error {
	a.initLogger()

	dst := ""
	if len(a.opts.args) != 0 {
		dst = a.opts.args[0]
	}

	if dst == "" {
		wd, err := env.Getwd()
		if err != nil {
			return fmt.Errorf("getting current working directory: %w", err)
		}

		dst = filepath.Join(wd, config.ConfigFileName())
	}

	info, err := os.Stat(dst)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("getting file info for destination path: %w", err)
		}
	}

	if info != nil {
		if !info.IsDir() {
			a.logger.Info("file already exists, skipping",
				"path", info.Name(), "size", info.Size(), "modified", info.ModTime(),
			)

			return nil
		}

		dst = filepath.Join(dst, config.ConfigFileName())
	}

	text := static.ConfigTemplate
	if a.opts.Plain {
		text = stripCfgComments(text)
	}

	cfgTmpl, err := template.New(config.DefaultConfigFile).Parse(text)
	if err != nil {
		return fmt.Errorf("parsing embedded configuration template: %w", err)
	}

	data := templateData{
		AppName: AppName,
		Version: Version(),
		Name:    cleanSessionName(filepath.Base(filepath.Dir(dst))),
		Time:    time.Now(),
		DocsURL: "https://github.com/michenriksen/tmpl",
	}

	cfgFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("creating configuration file: %w", err)
	}
	defer cfgFile.Close()

	if err := cfgTmpl.Execute(cfgFile, data); err != nil {
		return fmt.Errorf("writing configuration file: %w", err)
	}

	a.logger.Info("configuration file created", "path", cfgFile.Name())

	return nil
}

func cleanSessionName(name string) string {
	name = cleanSessNameRE.ReplaceAllString(strings.TrimSpace(name), "_")
	return strings.Trim(name, "._-")
}

func stripCfgComments(text string) string {
	lines := strings.Split(text, "\n")
	b := strings.Builder{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 || strings.HasPrefix(trimmed, "#") {
			continue
		}

		b.WriteString(line + "\n")
	}

	return strings.TrimSpace(b.String())
}

type templateData struct {
	AppName string
	Time    time.Time
	DocsURL string
	Name    string
	Version string
}
