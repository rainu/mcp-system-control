package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mcp-system-control/approval"
	"net/http"
	"os"
	"strings"

	"mcp-system-control/config/model"
	toolCommand "mcp-system-control/config/model/command"
	"mcp-system-control/expression"
	mcpCommand "mcp-system-control/mcp/server/builtin/tools/command"
	http2 "mcp-system-control/mcp/server/builtin/tools/http"

	"github.com/fatih/color"
	"github.com/goccy/go-yaml"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/olekukonko/tablewriter"
	"github.com/rainu/go-yacl"
)

type NoopLooger struct {
}

func (n NoopLooger) Print(message string) {
}

func (n NoopLooger) Trace(message string) {
}

func (n NoopLooger) Debug(message string) {
}

func (n NoopLooger) Info(message string) {
}

func (n NoopLooger) Warning(message string) {
}

func (n NoopLooger) Error(message string) {
}

func (n NoopLooger) Fatal(message string) {
}

func checkHelp(c *model.Config, config *yacl.Config) {
	if c.Help.Arg {
		printHelpArgs(os.Stdout, config)
		os.Exit(0)
	} else if c.Help.Env {
		printHelpEnv(os.Stdout)
		os.Exit(0)
	} else if c.Help.Yaml {
		printHelpConfig(os.Stdout)
		os.Exit(0)
	} else if c.Help.GenYaml {
		generateYamlSkeleton(os.Stdout)
		os.Exit(0)
	} else if c.Help.DumpYaml {
		dumpYaml(os.Stdout, c)
		os.Exit(0)
	} else if c.Help.Expr {
		printHelpExpression(os.Stdout)
		os.Exit(0)
	} else if c.Help.Tool {
		printHelpTool(os.Stdout)
		os.Exit(0)
	}
}

func printHelpArgs(output io.Writer, config *yacl.Config) {
	fmt.Fprintf(output, "Usage of %s:\n", os.Args[0])
	fmt.Fprint(output, config.HelpFlags(
		yacl.WithFilter(func(a yacl.FieldInfo) bool {
			p := a.Path()
			if strings.HasPrefix(p, "profile") {
				return p != "profile.[].description" && p != "profile.[].icon"
			}
			return false
		}),
		yacl.WithFlagDecorators(yacl.FlagDecorators{
			Short: func(s string) string {
				return color.New(color.Bold).Sprint(s)
			},
			LongKey: func(s string) string {
				return color.New(color.Bold).Sprint(s)
			},
			LongValue: func(s string) string {
				return color.New(color.Underline).Sprint(s)
			},
			Usage: func(s string) string {
				return color.New(color.Italic).Sprint(s)
			},
			DefaultValue: func(s string) string {
				return color.New(color.Faint).Sprint(s) + color.GreenString(" (default)")
			},
		}),
	))
}

func printHelpEnv(output io.Writer) {
	fmt.Fprintf(output, "All arguments can be within environment variables:\n")
	fmt.Fprintf(output, "Envronment variables with the prefix ")
	fmt.Fprintf(output, EnvironmentPrefix)
	fmt.Fprintf(output, " will be used. For example:\n")
	fmt.Fprintf(output, EnvironmentPrefix+"0=--mcp.name=MyMcpServer\n")
	fmt.Fprintf(output, EnvironmentPrefix+"1=--mcp.sse.bindAddress=:8080\n")
}

func generateYamlSkeleton(output io.Writer) {
	skeleton := &model.Config{}
	yacl.NewConfig(skeleton).ApplyDefaults()

	yaml.NewEncoder(output).Encode(skeleton)
}

func dumpYaml(output io.Writer, c *model.Config) {
	skeleton := &model.Config{}

	yaml.NewEncoder(output).Encode(skeleton)
}

func printHelpConfig(output io.Writer) {
	fmt.Fprintf(output, "Each available argument can be transformed into the corresponding yaml path. For example: '--mcp.sse.bindAddress=:8080'\n")
	s := ":8080"
	yaml.NewEncoder(output, yaml.Indent(2)).Encode(model.Config{
		MCP: model.MCP{
			SSE: model.MCPSSE{
				BindAddress: &s,
			},
		},
	})

	fmt.Fprintf(output, "\nYaml lookup file locations:\n")
	for _, location := range yamlLookupLocations() {
		fmt.Fprintf(output, "  - %s\n", location)
	}
}

func printHelpExpression(output io.Writer) {
	fmt.Fprintf(output, "The expression language is JavaScript. You can use the following variables and functions:\n")
	fmt.Fprintf(output, "\nFunctions:\n")
	fmt.Fprintf(output, "  - %s(...args): writes a message to the console.\n", expression.FuncNameLog)

	js := bytes.Buffer{}
	je := json.NewEncoder(&js)
	je.SetIndent("   ", "  ")
	je.Encode(mcpCommand.CommandDescriptor{
		Name:      "/path/to/command",
		Arguments: []string{"arg1", "...argN"},
		Environment: map[string]string{
			"ENV_VAR": "value",
		},
		AdditionalEnvironment: map[string]string{
			"ADDITIONAL_ENV_VAR": "value",
		},
		WorkingDirectory: "/path/to/working/dir",
	})
	fmt.Fprintf(output, "  - %s(%s): run a command.\n", expression.FuncNameRun, strings.TrimSpace(js.String()))

	js = bytes.Buffer{}
	je = json.NewEncoder(&js)
	je.SetIndent("   ", "  ")
	je.Encode(http2.CallDescriptor{
		Method: http.MethodPost,
		Url:    "https://example.com",
		Header: map[string]string{
			"Content-Type": "application/json",
		},
		StringBody: `{"msg": "hello world"}`,
	})
	fmt.Fprintf(output, "  - %s(%s): do a http call.\n", expression.FuncNameFetch, strings.TrimSpace(js.String()))
}

func printHelpTool(output io.Writer) {
	fmt.Fprintf(output, "Tool-Usage:"+
		"\nYou can define many functions that can be used by the MCP-Server."+
		"\nThe functions can be given by argument, Environment or config file."+
		"\nThe fields are more or less the same for all three methods:\n")

	fmt.Fprint(output, yacl.NewConfig(&toolCommand.FunctionDefinition{}).HelpFlags())

	exampleDefs := []toolCommand.FunctionDefinition{
		{
			Name:        "createFile",
			Description: "This function creates a file.",
			Parameters: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"path": map[string]any{
						"type":        "string",
						"description": "The path to the file.",
					},
				},
				Required: []string{"path"},
			},
			Command: "/usr/bin/touch",
			Environment: map[string]string{
				"USER":  "rainu",
				"SHELL": "/bin/bash",
			},
			WorkingDir: "/tmp",
			Approval:   "true",
		},
		{
			Name:        "echo",
			Description: "This function echoes a message.",
			Parameters: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"message": map[string]any{
						"type":        "string",
						"description": "The message to echo.",
					},
				},
				Required: []string{"message"},
			},
			AdditionalEnvironment: map[string]string{
				"MCP_SYSTEM_CONTROL_ARGS": "$@",
			},
			Command:  "/usr/bin/echo",
			Approval: "false",
		},
	}

	fmt.Fprintf(output, "\nJSON:\n")

	fdm := map[string]toolCommand.FunctionDefinition{}
	for _, def := range exampleDefs {
		jsonDef, _ := json.MarshalIndent(def, "", " ")
		fmt.Fprintf(output, "\n%s\n", jsonDef)

		fdm[def.Name] = def
	}

	fmt.Fprintf(output, "\nYAML:\n\n")
	ye := yaml.NewEncoder(output, yaml.Indent(2))
	ye.Encode(model.Config{Custom: fdm})

	fmt.Fprintf(output, "\nThe approval is always an js-expression. It will be evaluated each time the MCP-Server calls the function.\n")
	fmt.Fprintf(output, "If the expression returns true, the user must give the approval before the function will be executed.\n")
	fmt.Fprintf(output, "If the expression returns false, the user will NOT be asked for his approval.\n")
	fmt.Fprintf(output, "You can use the same variables and functions which are available in all other expressions (see --help-expression):\n")
	fmt.Fprintf(output, "The expression have access to the raw and parsed arguments from the MCP-Server and the function definition itself:\n")
	fmt.Fprintf(output, "  const %s = ", expression.VarNameContext)

	je := json.NewEncoder(output)
	je.SetIndent("  ", "  ")

	exampleDefs[0].Approval = `!` + expression.VarNameContext + `.args.path.startsWith('/tmp/')`
	je.Encode(approval.Variables{
		ToolDefinition: exampleDefs[0],
		RawArguments:   `{"path": "/tmp/file"}`,
		ParsedArguments: map[string]any{
			"path": "/tmp/file",
		},
	})

	fmt.Fprintf(output, "\nThe LLM will respond the arguments as JSON. You can use the following placeholders in the command:\n")
	fmt.Fprintf(output, "  - $@: all arguments (1:1 the JSON from the LLM)\n")
	fmt.Fprintf(output, "  - $<varName>: the value of <varName> in the LLM's JSON\n")
	fmt.Fprintf(output, "\nExamples:\n")

	table := tablewriter.NewWriter(output)
	table.SetBorder(false)
	table.SetHeader([]string{"Pattern", "LLM's JSON", "Result"})
	table.SetAutoWrapText(false)

	table.Append([]string{`/usr/bin/echo $@`, `{"message": "hello world"}`, `/usr/bin/echo {"message": "hello world"}`})
	table.Append([]string{`/usr/bin/echo $message`, `{"message": "hello world"}`, `/usr/bin/echo hello world`})
	table.Append([]string{`/usr/bin/echo "$message"`, `{"message": "hello world"}`, `/usr/bin/echo "hello world"`})
	table.Append([]string{`/usr/bin/echo "$message"`, `{}`, `/usr/bin/echo ""`})

	table.Render()

	fmt.Fprintf(output, "\nYou can also use these placeholder in (additional) environment and working directory variables.\n")

	fmt.Fprintf(output, "\nIt is also possible to define a JavaScript expression (file).\n")
	fmt.Fprintf(output, "You can use the same variables and functions which are available in all other expressions (see --help-expression):\n")
	fmt.Fprintf(output, "Additional variables:\n")
	fmt.Fprintf(output, "  const %s = ", expression.VarNameContext)

	je = json.NewEncoder(output)
	je.SetIndent("  ", "  ")
	je.Encode(toolCommand.Variables{
		FunctionDefinition: toolCommand.FunctionDefinition{
			Name:        "jsEcho",
			Description: "This function echoes a message.",
			Parameters: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"message": map[string]any{
						"type":        "string",
						"description": "The message to echo.",
					},
				},
				Required: []string{"message"},
			},
			CommandExpr: fmt.Sprintf(`"Echo: " + JSON.parse(%s.args).message`, expression.VarNameContext),
		},
		Arguments: `{"message": "hello world"}`,
	})

	fmt.Fprintf(output, "\nJavaScript command expression examples:")
	fmt.Fprintf(output, "\n\n  // parse llm's JSON, run the command and return its result\n")
	fmt.Fprintf(output, `
  const pa = JSON.parse(v.args)
  const cmdDescriptor = {
   "command": "echo",
   "arguments": ["Echo:", pa.message]
  }
  
  `+expression.FuncNameRun+`(cmdDescriptor)`)

	fmt.Fprintf(output, "\n\n  // catches possible execution error\n")
	fmt.Fprintf(output, `
  let result = ""
  try {
     result = `+expression.FuncNameRun+`({ "command": "__doesNotExists__" })
  } catch (e) {
     result = "Error: " + e
  }
  result
`)

}
