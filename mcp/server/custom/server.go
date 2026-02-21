package custom

import (
	"context"
	"encoding/json"
	"mcp-system-control/approval"

	"mcp-system-control/config/model/command"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewServer(version string, cfg map[string]command.FunctionDefinition, approvalRequester approval.Requester) *server.MCPServer {
	s := server.NewMCPServer(
		"mcp-system-control",
		version,
		server.WithToolCapabilities(false),
	)
	AddTools(s, cfg, approvalRequester)

	return s
}

func handlerFor(definition command.FunctionDefinition) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw, err := json.Marshal(request.Params.Arguments)
		if err != nil {
			return nil, err
		}

		rawResult, err := definition.CommandFn(ctx, string(raw))
		return mcp.NewToolResultText(string(rawResult)), err
	}
}

func AddTools(s *server.MCPServer, cfg map[string]command.FunctionDefinition, approvalRequester approval.Requester) {
	for name, definition := range cfg {
		t := mcp.Tool{
			Name:        name,
			Description: definition.Description,
			InputSchema: definition.Parameters,
		}
		s.AddTool(t, handlerFor(definition))
	}
}
