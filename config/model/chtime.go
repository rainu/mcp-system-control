package model

import (
	"mcp-system-control/config/model/approval"
)

type ChangeTimes struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`
}

func (c *ChangeTimes) SetDefaults() {
	if c.Approval == "" {
		c.Approval = approval.Never
	}
}
