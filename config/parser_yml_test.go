package config

import (
	"mcp-system-control/config/model"
	"mcp-system-control/config/model/command"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/go-yacl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_processYaml(t *testing.T) {
	c := &model.Config{}

	yamlContent := `
builtin:
  command-execution:
    disable: true
custom:
  test:
    description: This is a test function.
    parameters:
      type: object
      properties:
        arg1:
          type: string
          description: The first argument.
        arg2:
          type: number
          description: The second argument.
      required:
        - arg1
    command: doTest.sh
    approval: true
log-level: debug
`
	sr := strings.NewReader(yamlContent)
	config := yacl.NewConfig(c, yacl.WithAutoApplyDefaults(false))

	require.NoError(t, processYaml(config, sr))

	assert.Equal(t, &model.Config{
		DebugConfig: model.DebugConfig{
			LogLevel: "debug",
		},
		BuiltIns: model.BuiltIns{
			CommandExec: model.CommandExecution{
				Disable: true,
			},
		},
		Custom: map[string]command.FunctionDefinition{
			"test": {
				Description: "This is a test function.",
				Parameters: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]any{
						"arg1": map[string]any{
							"type":        "string",
							"description": "The first argument.",
						},
						"arg2": map[string]any{
							"type":        "number",
							"description": "The second argument.",
						},
					},
					Required: []string{"arg1"},
				},
				Command:  "doTest.sh",
				Approval: "true",
			},
		},
	}, c)
}
