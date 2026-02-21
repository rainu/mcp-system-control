package server

import (
	"mcp-system-control/approval"
	"mcp-system-control/config/model"
	"mcp-system-control/config/model/command"
	bServer "mcp-system-control/mcp/server/builtin"
	cServer "mcp-system-control/mcp/server/custom"

	"github.com/mark3labs/mcp-go/server"
)

func NewServer(name, version string, bConfig model.BuiltIns, cConfig map[string]command.FunctionDefinition, approvalRequester approval.Requester) *server.MCPServer {
	s := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(false),
	)
	bServer.AddTools(s, bConfig, approvalRequester)
	cServer.AddTools(s, cConfig, approvalRequester)

	return s
}
