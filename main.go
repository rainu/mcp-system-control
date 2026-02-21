package main

import (
	"fmt"
	"log/slog"
	"mcp-system-control/approval"
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

	ms := mcpServer.NewServer(
		cfg.MCP.Name,
		versionLine(),
		cfg.BuiltIns,
		cfg.Custom,
		approval.NewRequester(cfg.Approval),
	)

	var err error

	if cfg.MCP.SSE.BindAddress != nil {
		s := server.NewSSEServer(ms, cfg.MCP.SSE.Options()...)
		slog.Info(fmt.Sprintf("Starting SSE server on http://%s%s", *cfg.MCP.SSE.BindAddress, s.CompleteSsePath()))

		err = s.Start(*cfg.MCP.SSE.BindAddress)
	} else if cfg.MCP.Streamable.BindAddress != nil {
		s := server.NewStreamableHTTPServer(ms, cfg.MCP.Streamable.Options()...)
		slog.Info(fmt.Sprintf("Starting streamable server on http://%s%s", *cfg.MCP.Streamable.BindAddress, cfg.MCP.Streamable.EndpointPath))

		err = s.Start(*cfg.MCP.Streamable.BindAddress)
	} else {
		slog.Info("Starting stdio server")
		err = server.ServeStdio(ms, cfg.MCP.Stdio.Options()...)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}
