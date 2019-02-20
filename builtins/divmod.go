package builtins

import (
	. "github.com/prologic/monkey-lang/object"
)

// Divmod ...
func Divmod(args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}

	if a, ok := args[0].(*Integer); ok {
		if b, ok := args[1].(*Integer); ok {
			elements := make([]Object, 2)
			elements[0] = &Integer{Value: a.Value / b.Value}
			elements[1] = &Integer{Value: a.Value % b.Value}
			return &Array{Elements: elements}
		} else {
			return newError("expected argument #2 to `divmod` to be `int` got=%s", args[1].Type())
		}
	} else {
		return newError("expected argument #1 to `divmod` to be `int` got=%s", args[0].Type())
	}
}
