//go:build gen

//go:generate go run -tags=gen usage.go -f ../../docs/cli-usage.txt

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/michenriksen/tmpl/internal/cli"
)

func main() {
	writeOpt := flag.String("f", "", "write rendered template to file instead of stdout")
	flag.Parse()

	w := os.Stdout

	if *writeOpt != "" {
		var err error

		w, err = os.Create(*writeOpt)
		if err != nil {
			log.Fatal(fmt.Errorf("creating output file: %w", err))
		}
	}

	app, err := cli.NewApp(cli.WithOutputWriter(w))
	if err != nil {
		log.Fatal(fmt.Errorf("creating app: %w", err))
	}

	err = app.Run(context.Background(), "--help")
	if !errors.Is(err, cli.ErrHelp) {
		log.Fatal(fmt.Errorf("app did not return expected helper error: %w", err))
	}
}
