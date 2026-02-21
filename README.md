[![Go](https://github.com/rainu/mcp-system-control/actions/workflows/build.yml/badge.svg)](https://github.com/rainu/mcp-system-control/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/rainu/mcp-system-control/branch/main/graph/badge.svg)](https://codecov.io/gh/rainu/mcp-system-control)
[![Go Report Card](https://goreportcard.com/badge/github.com/rainu/mcp-system-control)](https://goreportcard.com/report/github.com/rainu/mcp-system-control)
[![Go Reference](https://pkg.go.dev/badge/github.com/rainu/mcp-system-control.svg)](https://pkg.go.dev/github.com/rainu/mcp-system-control)

# mcp-system-control

This is an MCP-Server which allows a LLM to control **the system**. 
It provides many tools which can be used by the LLM to control the system. 

For example:
- executing commands on the system
- file system access
- get system information

Each tool can be disabled. It is also possible to define your own custom tools which can be used by the LLM.
Each tool call can also be protected by a user approval. So before a tool is executed, the user will be asked if
they want to execute the tool call. This can be useful to prevent accidental tool calls which can be harmful for the system.

## Usage

The MCP-Server can be started in tree different modes:

### As stdio

```shell
mcp-system-control
```

### As SSE serer

```shell
mcp-system-control --mcp.sse.bindAddress=":8080"
```

### As Streamable server

```shell
mcp-system-control --mcp.streamable.bindAddress=":8080"
```

### Test

To test the server, you can use the [following command](https://modelcontextprotocol.io/docs/tools/inspector):

```shell
npx @modelcontextprotocol/inspector
```

## User approval

All tools can be protected by a user approval. If a tool is protected by a user approval, 
the user will be asked if they want to execute the tool call. For asking a user, system tools are used. The server will
try if one of these tools are available and use the first one which is available:
* [notify-send](https://man.archlinux.org/man/notify-send.1.en)
* [kdialog](https://github.com/KDE/kdialog)
* [zenity](https://github.com/ncruces/zenity)

If no system tool is available, the tool call will be rejected and an error will be returned to the LLM.

