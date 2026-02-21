package expression

import (
	"context"
	"mcp-system-control/mcp/server/builtin/tools/http"

	"github.com/dop251/goja"
)

const FuncNameFetch = "fetch"

func fetch(ctx context.Context, vm *goja.Runtime) func(http.CallDescriptor) *http.CallResult {
	return func(call http.CallDescriptor) *http.CallResult {
		r, err := call.Run(ctx, http.DefaultClient)
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}

		return r
	}
}
