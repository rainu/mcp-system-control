package approval

import (
	"context"
	"fmt"
	"os/exec"

	cfgModel "mcp-system-control/config/model/approval"

	"github.com/mark3labs/mcp-go/mcp"
)

type zenityRequester struct {
	cfg cfgModel.ZenityConfig
}

func newZenityRequester(cfg cfgModel.ZenityConfig) internalRequester {
	return &zenityRequester{cfg: cfg}
}

func (r *zenityRequester) WaitForApproval(ctx context.Context, request *mcp.CallToolRequest) (bool, error) {
	message := formatApprovalMessage(request)

	cmd := exec.CommandContext(ctx, "zenity",
		"--question",
		"--title="+r.cfg.Title,
		"--text="+message,
		fmt.Sprintf("--width=%d", r.cfg.Width),
		"--ok-label="+r.cfg.OkLabel,
		"--cancel-label="+r.cfg.CancelLabel,
	)

	err := cmd.Run()
	if err != nil {
		// Exit code 0 = OK/Approve, Exit code 1 = Cancel/Deny
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return false, nil // User denied
			}
		}
		return false, err
	}

	return true, nil // User approved
}

func (r *zenityRequester) IsAvailable() bool {
	return isCommandAvailable("zenity")
}
