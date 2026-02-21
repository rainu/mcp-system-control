package approval

import (
	"context"
	"encoding/json"
	"fmt"
	cfgModel "mcp-system-control/config/model/approval"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

type requester struct {
	cfg      cfgModel.Approval
	delegate internalRequester
}

func NewRequester(cfg cfgModel.Approval) Requester {
	result := requester{
		cfg: cfg,
	}

	switch cfg.Requester {
	case cfgModel.RequesterZenity:
		if isCommandAvailable("zenity") {
			result.delegate = newZenityRequester(cfg.Zenity)
		}
	case cfgModel.RequesterKDialog:
		if isCommandAvailable("kdialog") {
			result.delegate = newKDialogRequester(cfg.KDialog)
		}
	case cfgModel.RequesterNotifySend:
		if isCommandAvailable("notify-send") {
			result.delegate = newNotifySendRequester(cfg.NotifySend)
		}
	case cfgModel.RequesterCustom:
		r := newCustomRequester(cfg.Custom)
		if r.IsAvailable() {
			result.delegate = r
		}
	case cfgModel.RequesterAuto:
		fallthrough
	default:
		// Try requesters in order, only initializing when available
		if isCommandAvailable("notify-send") {
			result.delegate = newNotifySendRequester(cfg.NotifySend)
		} else if isCommandAvailable("zenity") {
			result.delegate = newZenityRequester(cfg.Zenity)
		} else if isCommandAvailable("kdialog") {
			result.delegate = newKDialogRequester(cfg.KDialog)
		} else {
			result.delegate = newCustomRequester(cfg.Custom)
		}
	}

	return &result
}

func (r *requester) WaitForApproval(ctx context.Context, request *mcp.CallToolRequest) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()

	if r.delegate != nil {
		return r.delegate.WaitForApproval(ctx, request)
	}
	return false, fmt.Errorf("unable to request approval to user")
}

// formatApprovalMessage formats the tool request into a human-readable message
func formatApprovalMessage(request *mcp.CallToolRequest) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Tool: %s\n\n", request.Params.Name))

	if request.Params.Arguments != nil {
		// Try to format arguments nicely
		if argsBytes, err := json.MarshalIndent(request.Params.Arguments, "", "  "); err == nil {
			sb.WriteString("Arguments:\n")
			sb.WriteString(string(argsBytes))
		} else {
			sb.WriteString(fmt.Sprintf("Arguments: %v", request.Params.Arguments))
		}
	}

	return sb.String()
}

// isCommandAvailable checks if a command is available in PATH
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
