package builtins

import (
	. "github.com/prologic/monkey-lang/object"
)

// Abs ...
func Abs(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if i, ok := args[0].(*Integer); ok {
		value := i.Value
		if value < 0 {
			value = value * -1
		}
		return &Integer{Value: value}
	}
	return newError("argument to `abs` not supported, got %s", args[0].Type())
}
