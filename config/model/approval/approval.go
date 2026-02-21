package approval

import "time"

type RequesterType string

const (
	RequesterAuto       RequesterType = "auto"
	RequesterZenity     RequesterType = "zenity"
	RequesterKDialog    RequesterType = "kdialog"
	RequesterNotifySend RequesterType = "notify-send"
	RequesterCustom     RequesterType = "custom"
)

type Approval struct {
	Timeout   time.Duration `yaml:"timeout,omitempty" usage:"Timeout for user"`
	Requester RequesterType `yaml:"requester,omitempty" usage:"Requester type to use (auto, zenity, kdialog, notify-send, custom)"`

	// Tool-specific configurations
	Zenity     ZenityConfig     `yaml:"zenity,omitempty" usage:"Zenity-specific: "`
	KDialog    KDialogConfig    `yaml:"kdialog,omitempty" usage:"KDialog-specific: "`
	NotifySend NotifySendConfig `yaml:"notify_send,omitempty" usage:"NotifySend-specific: "`
	Custom     CustomConfig     `yaml:"custom,omitempty" usage:"Custom script-based: "`
}

type ZenityConfig struct {
	Title       string `yaml:"title,omitempty" usage:"Title of the dialog window"`
	Width       int    `yaml:"width,omitempty" usage:"Width of the dialog window"`
	OkLabel     string `yaml:"ok_label,omitempty" usage:"Label for the OK button"`
	CancelLabel string `yaml:"cancel_label,omitempty" usage:"Label for the Cancel button"`
}

func (c *ZenityConfig) SetDefaults() {
	if c.Title == "" {
		c.Title = "MCP Tool Approval Required"
	}
	if c.Width == 0 {
		c.Width = 500
	}
	if c.OkLabel == "" {
		c.OkLabel = "Approve"
	}
	if c.CancelLabel == "" {
		c.CancelLabel = "Deny"
	}
}

type KDialogConfig struct {
	Title string `yaml:"title,omitempty" usage:"Title of the dialog window"`
}

func (c *KDialogConfig) SetDefaults() {
	if c.Title == "" {
		c.Title = "MCP Tool Approval Required"
	}
}

type NotifySendConfig struct {
	Urgency      string `yaml:"urgency,omitempty" usage:"Urgency level (low, normal, critical)"`
	Title        string `yaml:"title,omitempty" usage:"Title of the notification"`
	ApproveLabel string `yaml:"approve_label,omitempty" usage:"Label for the Approve action"`
	DenyLabel    string `yaml:"deny_label,omitempty" usage:"Label for the Deny action"`
}

func (c *NotifySendConfig) SetDefaults() {
	if c.Urgency == "" {
		c.Urgency = "critical"
	}
	if c.Title == "" {
		c.Title = "MCP Tool Approval Required"
	}
	if c.ApproveLabel == "" {
		c.ApproveLabel = "Approve"
	}
	if c.DenyLabel == "" {
		c.DenyLabel = "Deny"
	}
}

type CustomConfig struct {
	Script string   `yaml:"script,omitempty" usage:"Path to the custom script to execute for approval"`
	Args   []string `yaml:"args,omitempty" usage:"Additional arguments to pass to the script"`
}

func (c *Approval) SetDefaults() {
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}
	if c.Requester == "" {
		c.Requester = RequesterAuto
	}
}
