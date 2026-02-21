package approval

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"

	cfgModel "mcp-system-control/config/model/approval"

	"github.com/mark3labs/mcp-go/mcp"
)

type customRequester struct {
	cfg cfgModel.CustomConfig
}

func newCustomRequester(cfg cfgModel.CustomConfig) internalRequester {
	return &customRequester{cfg: cfg}
}

func (r *customRequester) WaitForApproval(ctx context.Context, request *mcp.CallToolRequest) (bool, error) {
	// Serialize the request to JSON
	requestJSON, err := json.Marshal(request.Params)
	if err != nil {
		return false, err
	}

	// Build the command with script path and args
	args := append(r.cfg.Args, string(requestJSON))
	cmd := exec.CommandContext(ctx, r.cfg.Script, args...)

	// Run the script
	err = cmd.Run()
	if err != nil {
		// Check if it's an exit error
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Exit code 0 = Approved, anything else = Denied
			return exitErr.ExitCode() == 0, nil
		}
		return false, err
	}

	// Exit code 0 = Approved
	return true, nil
}

func (r *customRequester) IsAvailable() bool {
	// Check if the script exists and is executable
	if r.cfg.Script == "" {
		return false
	}

	info, err := os.Stat(r.cfg.Script)
	if err != nil {
		return false
	}

	// Check if it's a regular file and executable
	mode := info.Mode()
	if !mode.IsRegular() {
		return false
	}

	// Check if file is executable (has any execute bit set)
	return mode&0111 != 0
}
