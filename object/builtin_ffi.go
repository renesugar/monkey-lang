package object

import (
	"fmt"
	"plugin"
)

// FFI ...
func FFI(args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}

	arg, ok := args[0].(*String)
	if !ok {
		return newError("argument #1 to `ffi` expected to be `str` got=%T", args[0].Type())
	}
	name := arg.Value

	arg, ok = args[1].(*String)
	if !ok {
		return newError("argument #2 to `ffi` expected to be `str` got=%T", args[0].Type())
	}
	symbol := arg.Value

	p, err := plugin.Open(fmt.Sprintf("%s.so", name))
	if err != nil {
		return newError("error loading plugin: %s", err)
	}

	v, err := p.Lookup(symbol)
	if err != nil {
		return newError("error finding symbol: %s", err)
	}

	return &Builtin{Name: symbol, Fn: v.(func(...Object) Object)}
}
