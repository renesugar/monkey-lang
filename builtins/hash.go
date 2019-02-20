package builtins

import (
	. "github.com/prologic/monkey-lang/object"
)

// HashOf ...
func HashOf(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if hash, ok := args[0].(Hashable); ok {
		return &Integer{Value: int64(hash.HashKey().Value)}
	}
	return newError("argument #1 to `hash()` is not hashable: %s", args[0].Inspect())
}
