//go:build gen

//go:generate go run -tags=gen readme.go -f ../../README.md

// Package gen contains code for generating documentation and other files.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const (
	docsURL        = "https://michenriksen.com/tmpl"
	projectPackage = "github.com/michenriksen/tmpl"
	projectURL     = "https://github.com/michenriksen/tmpl"
)

func main() {
	writeOpt := flag.String("f", "", "write rendered template to file instead of stdout")
	flag.Parse()

	tmplPath := filepath.Join("docs", "README.md.tmpl")

	tmpl, err := template.New("README.md.tmpl").Funcs(helpers).ParseFiles(tmplPath)
	if err != nil {
		log.Fatal(fmt.Errorf("parsing template: %w", err))
	}

	out := os.Stdout
	if *writeOpt != "" {
		out, err = os.Create(*writeOpt)
		if err != nil {
			log.Fatal(fmt.Errorf("creating output file: %w", err))
		}
	}

	err = tmpl.Funcs(helpers).Execute(out, map[string]any{
		"DocsURL":        docsURL,
		"ProjectPackage": projectPackage,
		"ProjectURL":     projectURL,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("rendering template: %w", err))
	}
}

var helpers = template.FuncMap{
	"file": func(name string) (string, error) {
		data, err := os.ReadFile(name)
		if err != nil {
			return "", fmt.Errorf("reading file: %w", err)
		}

		return string(data), nil
	},
}
