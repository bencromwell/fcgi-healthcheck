// Package cmd provides the CLI for sshush.
package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/bencromwell/fcgi-healthcheck/fcgi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func must(err error) {
	if err != nil {
		slog.Error("fgci-healthcheck", "error", err)
		os.Exit(1)
	}
}

// NewRootCommand creates a new root command.
func NewRootCommand(version, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "fcgi-healthcheck",
		Short:         "fcgi-healthcheck",
		Version:       fmt.Sprintf("%s (%s)", version, commit),
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, _ []string) error {
			params := &fcgi.Input{
				Config: &fcgi.Config{
					Network:      viper.GetString("network"),
					Addr:         viper.GetString("addr"),
					DocumentRoot: viper.GetString("document-root"),
					RemoteAddr:   viper.GetString("remote-addr"),
					ServerAddr:   viper.GetString("server-addr"),
					ServerPort:   viper.GetString("server-port"),
					ServerName:   viper.GetString("server-name"),
				},
				Request: &fcgi.Request{
					Script: viper.GetString("script"),
					URI:    viper.GetString("uri"),
					Method: viper.GetString("method"),
					Query:  viper.GetString("query"),
					Body:   viper.GetString("body"),
				},
				Verbose: viper.GetBool("verbose"),
			}

			err := fcgi.Run(params)

			return err
		},
	}

	cmd.PersistentFlags().BoolP("verbose", "V", false, "verbose output")

	cmd.PersistentFlags().String(
		"script",
		"",
		"Path to PHP script on disk (as seen by PHP-FPM)",
	)

	cmd.PersistentFlags().String(
		"uri",
		"",
		"REQUEST_URI",
	)

	cmd.PersistentFlags().String(
		"method",
		"GET",
		"HTTP method",
	)

	cmd.PersistentFlags().String(
		"query",
		"",
		"QUERY_STRING (e.g. a=1&b=2)",
	)

	cmd.PersistentFlags().String(
		"body",
		"",
		"Request body (for POST, etc.)",
	)

	// Request flags are bound to Viper.
	must(viper.BindPFlag("verbose", cmd.PersistentFlags().Lookup("verbose")))
	must(viper.BindPFlag("script", cmd.PersistentFlags().Lookup("script")))
	must(viper.BindPFlag("uri", cmd.PersistentFlags().Lookup("uri")))
	must(viper.BindPFlag("method", cmd.PersistentFlags().Lookup("method")))
	must(viper.BindPFlag("query", cmd.PersistentFlags().Lookup("query")))
	must(viper.BindPFlag("body", cmd.PersistentFlags().Lookup("body")))

	// Config-only values have defaults set direct in Viper.
	// network can be tcp or unix.
	viper.SetDefault("network", "tcp")
	// FastCGI address is host:port for tcp, path for unix.
	viper.SetDefault("addr", "127.0.0.1:9000")
	viper.SetDefault("document-root", "/var/www/html")
	viper.SetDefault("remote-addr", "127.0.0.1")
	viper.SetDefault("server-addr", "127.0.0.1")
	viper.SetDefault("remote-port", "80")
	viper.SetDefault("server-port", "80")
	viper.SetDefault("server-name", "localhost")

	// Config files.
	viper.SetConfigName("fgci")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Envs.
	viper.SetEnvPrefix("FCGI")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := viper.ReadInConfig(); err != nil {
		// Config file is optional.
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			must(fmt.Errorf("reading config: %w", err))
		}
	}

	return cmd
}
