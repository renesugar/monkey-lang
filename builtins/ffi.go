package builtins

import (
	"fmt"
	"plugin"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// FFI ...
func FFI(args ...object.Object) object.Object {
	if err := typing.Check(
		"ffi", args,
		typing.ExactArgs(2),
		typing.WithTypes(object.STRING, object.STRING),
	); err != nil {
		return newError(err.Error())
	}

	name := args[0].(*object.String).Value
	symbol := args[1].(*object.String).Value

	p, err := plugin.Open(fmt.Sprintf("%s.so", name))
	if err != nil {
		return newError("error loading plugin: %s", err)
	}

	v, err := p.Lookup(symbol)
	if err != nil {
		return newError("error finding symbol: %s", err)
	}

	return &object.Builtin{Name: symbol, Fn: v.(object.BuiltinFunction)}
}
