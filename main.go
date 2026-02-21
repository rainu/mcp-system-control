package main

import (
	"fmt"
	"log/slog"
	"mcp-system-control/config"
	mcpServer "mcp-system-control/mcp/server"
	"os"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	cfg := config.Parse(os.Args[1:], os.Environ())
	if cfg.Version {
		fmt.Fprintln(os.Stderr, versionLine())
		os.Exit(0)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	slog.SetLogLoggerLevel(*cfg.DebugConfig.LogLevelParsed)

	ms := mcpServer.NewServer(versionLine(), cfg.BuiltIns, cfg.Custom)

	var err error
	if cfg.HttpAddress != "" {
		slog.Info(fmt.Sprintf("Starting streamable http server on http://%s/mcp", cfg.HttpAddress))
		err = server.NewStreamableHTTPServer(ms).Start(cfg.HttpAddress)
	} else {
		err = server.ServeStdio(ms)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}
