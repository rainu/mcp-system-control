package system

import (
	"context"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

var SystemTimeTool = mcp.NewTool("getSystemTime",
	mcp.WithDescription("Get the current system time."),
)
var SystemTimeToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText(time.Now().String()), nil
}
