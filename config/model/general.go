package model

import (
	"fmt"
	"mcp-system-control/config/model/command"
)

type Config struct {
	ConfigFile `yaml:",inline"`

	DebugConfig DebugConfig `yaml:",inline,omitempty"`

	MCP MCP `yaml:"mcp,omitempty" usage:"MCP server configuration: "`

	BuiltIns BuiltIns                              `yaml:"builtin,omitempty" usage:"Built-in tool "`
	Custom   map[string]command.FunctionDefinition `yaml:"custom,omitempty" usage:"Custom tool definition "`

	Version bool `yaml:"version,omitempty" short:"v" usage:"Show the version"`

	Help Help `yaml:",inline,omitempty"`
}

type ConfigFile struct {
	Path string `yaml:"config-file,omitempty" short:"c" usage:"Path to the configuration yaml file"`
}

func (c *Config) Validate() error {
	if ve := c.DebugConfig.Validate(); ve != nil {
		return ve
	}

	for cmd, definition := range c.Custom {
		definition.Name = cmd

		if definition.Parameters.Type == "" && len(definition.Parameters.Properties) == 0 {
			definition.Parameters.Type = "object"                   // Default to object if no type is set
			definition.Parameters.Properties = make(map[string]any) // Ensure Properties is initialized
		}

		if definition.CommandExpr != "" {
			if ve := command.Expression(definition.CommandExpr).Validate(); ve != nil {
				return ve
			}
			definition.CommandFn = command.Expression(definition.CommandExpr).CommandFn(definition)
		} else if definition.Command != "" {
			if ve := command.Command(definition.Command).Validate(); ve != nil {
				return ve
			}
			definition.CommandFn = command.Command(definition.Command).CommandFn(definition)
		} else {
			return fmt.Errorf("Command for tool '%s' is missing", cmd)
		}

		// definition is only a local copy, so we need to set it back
		c.Custom[cmd] = definition
	}

	return nil
}
