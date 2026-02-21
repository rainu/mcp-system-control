package server

import (
	"context"
	"log/slog"
	"mcp-system-control/approval"
	"mcp-system-control/config/model"
	"mcp-system-control/config/model/command"
	bServer "mcp-system-control/mcp/server/builtin"
	cServer "mcp-system-control/mcp/server/custom"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewServer(name, version string, bConfig model.BuiltIns, cConfig map[string]command.FunctionDefinition, approvalRequester approval.Requester) *server.MCPServer {
	s := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(false),
		server.WithHooks(&server.Hooks{
			OnBeforeAny: []server.BeforeAnyHookFunc{
				func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
					if method != mcp.MethodToolsCall {
						slog.Info("Action",
							slog.Any("id", id),
							slog.Any("method", method),
						)
					}
				},
			},
			OnBeforeCallTool: []server.OnBeforeCallToolFunc{
				func(ctx context.Context, id any, message *mcp.CallToolRequest) {
					slog.Info("Tool call",
						slog.Any("id", id),
						slog.String("tool", message.Params.Name),
					)
				},
			},
		}),
	)
	bServer.AddTools(s, bConfig, approvalRequester)
	cServer.AddTools(s, cConfig, approvalRequester)

	return s
}
