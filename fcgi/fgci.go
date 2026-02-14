package fcgi

import (
	"fmt"
	"io"
	"os"
	"strings"

	fcgi "github.com/tomasen/fcgi_client"
)

type (
	Config struct {
		Network      string
		Addr         string
		DocumentRoot string
		RemoteAddr   string
		RemotePort   string
		ServerAddr   string
		ServerPort   string
		ServerName   string
	}

	Request struct {
		Script string
		URI    string
		Method string
		Query  string
		Body   string
	}

	Input struct {
		Request *Request
		Config  *Config
		Verbose bool
	}
)

func (i *Input) toParams() map[string]string {
	return map[string]string{
		"REQUEST_METHOD":  i.Request.Method,
		"REQUEST_URI":     i.Request.URI,
		"SCRIPT_NAME":     i.Request.URI,
		"SCRIPT_FILENAME": i.Request.Script,
		"QUERY_STRING":    i.Request.Query,

		"DOCUMENT_ROOT": i.Config.DocumentRoot,
		"REMOTE_ADDR":   i.Config.RemoteAddr,
		"REMOTE_PORT":   i.Config.RemotePort,
		"SERVER_ADDR":   i.Config.ServerAddr,
		"SERVER_PORT":   i.Config.ServerPort,
		"SERVER_NAME":   i.Config.ServerName,
	}
}

func Run(params *Input) error {
	if len(params.Request.URI) == 0 && len(params.Request.Script) == 0 {
		return fmt.Errorf("you must supply either uri or script")
	}

	client, err := fcgi.Dial(params.Config.Network, params.Config.Addr)
	if err != nil {
		return fmt.Errorf("dial php-fpm: %w", err)
	}
	defer client.Close()

	resp, err := client.Request(
		params.toParams(),
		strings.NewReader(params.Request.Body),
	)
	if err != nil {
		return fmt.Errorf("fcgi request failed: %w", err)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("err: %w", err)
	}

	if resp.StatusCode == 403 {
		return fmt.Errorf("forbidden: %v", strings.TrimRight(string(content), "\n"))
	}

	if params.Verbose {
		fmt.Fprintln(os.Stdout, string(content))
	}

	return nil
}
