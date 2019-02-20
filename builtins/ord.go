package builtins

import (
	. "github.com/prologic/monkey-lang/object"
)

// Ord ...
func Ord(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if s, ok := args[0].(*String); ok {
		if len(s.Value) == 1 {
			return &Integer{Value: int64(s.Value[0])}
		}
		return newError("`ord()` expected a character but got string of length %d", len(s.Value))
	}
	return newError("argument to `ord` not supported, got %s", args[0].Type())
}
