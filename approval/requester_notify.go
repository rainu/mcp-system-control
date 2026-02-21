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
		"-A", r.cfg.ApproveLabel,
		"-A", r.cfg.DenyLabel,
	}

	args = append(args, r.cfg.Title, message)

	cmd := exec.CommandContext(ctx, "notify-send", args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return false, err
	}

	// notify-send returns the action index when clicked
	// 0 = first action (Approve), 1 = second action (Deny)
	output := strings.TrimSpace(out.String())

	return output == "0", nil
}

func (r *notifySendRequester) IsAvailable() bool {
	return isCommandAvailable("notify-send")
}
