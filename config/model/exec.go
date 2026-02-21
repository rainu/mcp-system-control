package model

import (
	"mcp-system-control/config/model/approval"
)

type CommandExecution struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Needs no user approval to be executed"`
}

func (c *CommandExecution) SetDefaults() {
	if c.Approval == "" {
		c.Approval = approval.Always
	}
}
