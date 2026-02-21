package approval

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	cfgModel "mcp-system-control/config/model/approval"

	"github.com/mark3labs/mcp-go/mcp"
)

type notifySendRequester struct {
	cfg cfgModel.NotifySendConfig
}

func newNotifySendRequester(cfg cfgModel.NotifySendConfig) internalRequester {
	return &notifySendRequester{cfg: cfg}
}

func (r *notifySendRequester) WaitForApproval(ctx context.Context, request *mcp.CallToolRequest) (bool, error) {
	message := formatApprovalMessage(request)
	args := []string{
		"-u", r.cfg.Urgency,
		"-w", // Wait for action (blocks until notification is closed)
		"-A", r.cfg.DenyLabel,
		"-A", r.cfg.ApproveLabel,
	}

	args = append(args, r.cfg.Title, message)

	cmd := exec.CommandContext(ctx, "notify-send", args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		// If the command fails, treat it as denial
		return false, err
	}

	// notify-send returns the action index when clicked
	// 0 = first action (Deny), 1 = second action (Approve)
	// empty string or no output = notification closed without action (timeout/dismiss)
	output := strings.TrimSpace(out.String())

	// If no output, the notification was closed without clicking an action
	// This happens when the notification times out or is dismissed
	if output == "" {
		return false, nil // Treat timeout/dismiss as denial
	}

	return output == "1", nil
}

func (r *notifySendRequester) IsAvailable() bool {
	return isCommandAvailable("notify-send")
}
