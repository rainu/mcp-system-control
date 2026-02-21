package approval

import (
	"context"
	"fmt"
	"log/slog"
	cfgModel "mcp-system-control/config/model/approval"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
)

type requester struct {
	cfg      cfgModel.Approval
	delegate internalRequester
}

func NewRequester(cfg cfgModel.Approval) Requester {
	// Initialize language for approval messages
	if err := SetLanguage(cfg.Language); err != nil {
		slog.Warn("Failed to set approval message language, using auto-detect", "error", err)
	}

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

// isCommandAvailable checks if a command is available in PATH
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
