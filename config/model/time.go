package model

import (
	"mcp-system-control/config/model/approval"
)

type SystemTime struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`
}

func (c *SystemTime) SetDefaults() {
	if c.Approval == "" {
		c.Approval = approval.Never
	}
}
