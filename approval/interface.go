package approval

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type Requester interface {
	WaitForApproval(ctx context.Context, request *mcp.CallToolRequest) (bool, error)
}

type internalRequester interface {
	Requester
	IsAvailable() bool
}
