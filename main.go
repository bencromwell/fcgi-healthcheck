package main

import (
	"log/slog"
	"os"

	"github.com/bencromwell/fcgi-healthcheck/cmd"
)

//nolint:gochecknoglobals // These variables are set using ldflags.
var (
	version = "0.0.0-dev"
	commit  = ""
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	rootCmd := cmd.NewRootCommand(version, commit)

	err := rootCmd.Execute()
	if err != nil {
		// Cobra has already printed the error.
		os.Exit(1)
	}
}
