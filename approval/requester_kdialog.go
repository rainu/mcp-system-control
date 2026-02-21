package approval

import (
	"context"
	"os/exec"

	cfgModel "mcp-system-control/config/model/approval"

	"github.com/mark3labs/mcp-go/mcp"
)

type kdialogRequester struct {
	cfg cfgModel.KDialogConfig
}

func newKDialogRequester(cfg cfgModel.KDialogConfig) internalRequester {
	return &kdialogRequester{cfg: cfg}
}

func (r *kdialogRequester) WaitForApproval(ctx context.Context, request *mcp.CallToolRequest) (bool, error) {
	message := formatApprovalMessage(request)

	cmd := exec.CommandContext(ctx, "kdialog",
		"--yesno", message,
		"--title", r.cfg.Title,
	)

	err := cmd.Run()
	if err != nil {
		// Exit code 0 = Yes/Approve, Exit code 1 = No/Deny
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return false, nil // User denied
			}
		}
		return false, err
	}

	return true, nil // User approved
}

func (r *kdialogRequester) IsAvailable() bool {
	return isCommandAvailable("kdialog")
}
