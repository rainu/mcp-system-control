package model

import (
	"time"

	"github.com/mark3labs/mcp-go/server"
)

type MCP struct {
	Name string `yaml:"name,omitempty" usage:"Name of the MCP server"`

	SSE        MCPSSE        `yaml:"sse,omitempty" usage:"[SSE] "`
	Streamable MCPStreamable `yaml:"streamable,omitempty" usage:"[Streamable] "`
	Stdio      MCPStdio      `yaml:"stdio,omitempty" usage:"[Stdio] "`
}

func (c *MCP) SetDefaults() {
	if c.Name == "" {
		c.Name = "mcp-system-control"
	}
}

type MCPSSE struct {
	BindAddress                  *string        `yaml:"bindAddress,omitempty" usage:"Bind address for SSE server"`
	BasePath                     *string        `yaml:"basePath,omitempty" usage:"Base path for SSE server"`
	SSEEndpoint                  *string        `yaml:"sseEndpoint,omitempty" usage:"SSE endpoint path"`
	MessageEndpoint              *string        `yaml:"messageEndpoint,omitempty" usage:"Message endpoint path"`
	AppendQueryToMessageEndpoint *bool          `yaml:"appendQueryToMessageEndpoint,omitempty" usage:"Append query parameters to message endpoint"`
	KeepAlive                    *bool          `yaml:"keepAlive,omitempty" usage:"Enable SSE keep-alive"`
	KeepAliveInterval            *time.Duration `yaml:"keepAliveInterval,omitempty" usage:"SSE keep-alive interval"`
}

func (c *MCPSSE) Options() []server.SSEOption {
	var opts []server.SSEOption

	if c.BasePath != nil {
		opts = append(opts, server.WithStaticBasePath(*c.BasePath))
	}
	if c.SSEEndpoint != nil {
		opts = append(opts, server.WithSSEEndpoint(*c.SSEEndpoint))
	}
	if c.MessageEndpoint != nil {
		opts = append(opts, server.WithMessageEndpoint(*c.MessageEndpoint))
	}
	if c.AppendQueryToMessageEndpoint != nil && *c.AppendQueryToMessageEndpoint {
		opts = append(opts, server.WithAppendQueryToMessageEndpoint())
	}
	if c.KeepAlive != nil {
		opts = append(opts, server.WithKeepAlive(*c.KeepAlive))
	}
	if c.KeepAliveInterval != nil {
		opts = append(opts, server.WithKeepAliveInterval(*c.KeepAliveInterval))
	}

	return opts
}

type MCPStreamable struct {
	BindAddress       *string        `yaml:"bindAddress,omitempty" usage:"Bind address for SSE server"`
	EndpointPath      string         `yaml:"endpointPath,omitempty" usage:"HTTP endpoint path"`
	DisableStreaming  *bool          `yaml:"disableStreaming,omitempty" usage:"Disable streaming mode"`
	HeartbeatInterval *time.Duration `yaml:"heartbeatInterval,omitempty" usage:"Heartbeat interval for streaming"`
	Stateless         *bool          `yaml:"stateless,omitempty" usage:"Run in stateless mode"`
	Stateful          *bool          `yaml:"stateful,omitempty" usage:"Run in stateful mode"`
	TLSCertFile       *string        `yaml:"tlsCertFile,omitempty" usage:"TLS certificate file path"`
	TLSKeyFile        *string        `yaml:"tlsKeyFile,omitempty" usage:"TLS key file path"`
}

func (c *MCPStreamable) SetDefaults() {
	if c.EndpointPath == "" {
		c.EndpointPath = "/mcp"
	}
}

func (c *MCPStreamable) Options() []server.StreamableHTTPOption {
	var opts []server.StreamableHTTPOption

	if c.EndpointPath != "" {
		opts = append(opts, server.WithEndpointPath(c.EndpointPath))
	}
	if c.DisableStreaming != nil {
		opts = append(opts, server.WithDisableStreaming(*c.DisableStreaming))
	}
	if c.HeartbeatInterval != nil {
		opts = append(opts, server.WithHeartbeatInterval(*c.HeartbeatInterval))
	}
	if c.Stateless != nil && *c.Stateless {
		opts = append(opts, server.WithStateLess(true))
	}
	if c.Stateful != nil && *c.Stateful {
		opts = append(opts, server.WithStateful(true))
	}
	if c.TLSCertFile != nil && c.TLSKeyFile != nil {
		opts = append(opts, server.WithTLSCert(*c.TLSCertFile, *c.TLSKeyFile))
	}

	return opts
}

type MCPStdio struct {
	QueueSize      *int `yaml:"queueSize,omitempty" usage:"Queue size for stdio communication"`
	WorkerPoolSize *int `yaml:"workerPoolSize,omitempty" usage:"Worker pool size for stdio processing"`
}

func (c *MCPStdio) Options() []server.StdioOption {
	var opts []server.StdioOption

	if c.QueueSize != nil {
		opts = append(opts, server.WithQueueSize(*c.QueueSize))
	}
	if c.WorkerPoolSize != nil {
		opts = append(opts, server.WithWorkerPoolSize(*c.WorkerPoolSize))
	}

	return opts
}
