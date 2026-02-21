package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mcp-system-control/approval"
	"mcp-system-control/config/model"
	"mcp-system-control/mcp/server/builtin/tools/command"
	"mcp-system-control/mcp/server/builtin/tools/file"
	"mcp-system-control/mcp/server/builtin/tools/system"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func AddTools(s *server.MCPServer, cfg model.BuiltIns, approvalRequester approval.Requester) {
	addTool := func(tool mcp.Tool, handler server.ToolHandlerFunc) {
		as := approval.Approval(cfg.GetApprovalFor(tool.Name))

		if as == approval.Never {
			s.AddTool(tool, handler)
		} else if as == approval.Always {
			s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				approved, err := approvalRequester.WaitForApproval(ctx, &request)
				if err != nil {
					return nil, fmt.Errorf("error while waiting for approval: %w", err)
				}
				if !approved {
					return nil, fmt.Errorf("tool call not approved")
				}

				return handler(ctx, request)
			})
		} else {
			s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				argsAsJson, err := json.Marshal(request.Params.Arguments)
				if err != nil {
					slog.Error("Failed to marshal arguments", "error", err)
					return nil, fmt.Errorf("failed to marshal arguments")
				}
				if as.NeedsApproval(ctx, string(argsAsJson), nil) {
					approved, err := approvalRequester.WaitForApproval(ctx, &request)
					if err != nil {
						return nil, fmt.Errorf("error while waiting for approval: %w", err)
					}
					if !approved {
						return nil, fmt.Errorf("tool call not approved")
					}
				}
				return handler(ctx, request)
			})
		}
	}

	if !cfg.SystemTime.Disable {
		addTool(system.SystemTimeTool, system.SystemTimeToolHandler)
	}
	if !cfg.SystemInfo.Disable {
		addTool(system.SystemInfoTool, system.SystemInfoToolHandler)
	}
	if !cfg.Environment.Disable {
		addTool(system.EnvironmentTool, system.EnvironmentToolHandler)
	}

	if !cfg.ChangeMode.Disable {
		addTool(file.ChangeModeTool, file.ChangeModeToolHandler)
	}
	if !cfg.ChangeOwner.Disable {
		addTool(file.ChangeOwnerTool, file.ChangeOwnerToolHandler)
	}
	if !cfg.ChangeTimes.Disable {
		addTool(file.ChangeTimesTool, file.ChangeTimesToolHandler)
	}

	if !cfg.DirectoryCreation.Disable {
		addTool(file.DirectoryCreationTool, file.DirectoryCreationToolHandler)
	}
	if !cfg.DirectoryDeletion.Disable {
		addTool(file.DirectoryDeletionTool, file.DirectoryDeletionToolHandler)
	}
	if !cfg.DirectoryTempCreation.Disable {
		addTool(file.DirectoryTempCreationTool, file.DirectoryTempCreationToolHandler)
	}

	if !cfg.FileAppending.Disable {
		addTool(file.FileAppendingTool, file.FileAppendingToolHandler)
	}
	if !cfg.FileCreation.Disable {
		addTool(file.FileCreationTool, file.FileCreationToolHandler)
	}
	if !cfg.FileDeletion.Disable {
		addTool(file.FileDeletionTool, file.FileDeletionToolHandler)
	}
	if !cfg.FileReading.Disable {
		addTool(file.FileReadingTool, file.FileReadingToolHandler)
	}
	if !cfg.FileTempCreation.Disable {
		addTool(file.FileTempCreationTool, file.FileTempCreationToolHandler)
	}
	if !cfg.Stats.Disable {
		addTool(file.StatsTool, file.StatsToolHandler)
	}

	if !cfg.CommandExec.Disable {
		addTool(command.CommandExecutionTool, command.CommandExecutionToolHandler)
	}
}
