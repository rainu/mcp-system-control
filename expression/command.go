package expression

import (
	"context"
	"mcp-system-control/mcp/server/builtin/tools/command"

	"github.com/dop251/goja"
)

const FuncNameRun = "run"

func run(ctx context.Context, vm *goja.Runtime) func(command.CommandDescriptor) string {
	return func(cmd command.CommandDescriptor) string {
		r, err := cmd.Run(ctx)
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}

		return string(r)
	}
}
