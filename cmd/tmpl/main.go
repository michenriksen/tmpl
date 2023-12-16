package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/michenriksen/tmpl/internal/cli"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	exitChan := make(chan int, 1)

	signalHandler(cancel, exitChan)

	app, err := cli.NewApp()
	if err != nil {
		panic(err)
	}

	if err := app.Run(ctx, os.Args[1:]...); err != nil {
		if errors.Is(err, cli.ErrHelp) || errors.Is(err, cli.ErrVersion) {
			os.Exit(0)
		}

		// Return exit code 2 when a configuration file is invalid.
		if errors.Is(err, cli.ErrInvalidConfig) {
			os.Exit(2)
		}

		if !errors.Is(ctx.Err(), context.Canceled) {
			os.Exit(1)
		}

		exitCode := <-exitChan
		os.Exit(exitCode)
	}
}

func signalHandler(cancel context.CancelFunc, exitChan chan<- int) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range sigchan {
			cancel()

			fmt.Fprintf(os.Stderr, "received signal: %s; exiting...\n", sig)
			time.Sleep(1 * time.Second)

			exitChan <- 1
		}
	}()
}
