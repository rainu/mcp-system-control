package system

import (
	"os"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTool_Env(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(EnvironmentTool, EnvironmentToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = EnvironmentTool.Name

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	text := res.Content[0].(mcp.TextContent).Text

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)

		assert.Contains(t, text, pair[0])
		assert.Contains(t, text, pair[1])
	}
}

func getTestClient(t *testing.T, serverConf func(s *server.MCPServer)) *client.Client {
	s := server.NewMCPServer(
		"mcp-system-control",
		"test-version",
		server.WithToolCapabilities(false),
	)
	serverConf(s)

	c := client.NewClient(transport.NewInProcessTransport(s))

	_, err := c.Initialize(t.Context(), mcp.InitializeRequest{})
	require.NoError(t, err)

	return c
}
