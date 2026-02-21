package server

import (
	"mcp-system-control/config/model"
	"mcp-system-control/config/model/command"
	bServer "mcp-system-control/mcp/server/builtin"
	cServer "mcp-system-control/mcp/server/custom"

	"github.com/mark3labs/mcp-go/server"
)

func NewServer(version string, bConfig model.BuiltIns, cConfig map[string]command.FunctionDefinition) *server.MCPServer {
	s := server.NewMCPServer(
		"ask-mai",
		version,
		server.WithToolCapabilities(false),
	)
	bServer.AddTools(s, bConfig)
	cServer.AddTools(s, cConfig)

	return s
}
