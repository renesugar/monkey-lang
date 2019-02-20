package builtins

import (
	. "github.com/prologic/monkey-lang/object"
)

func pow(x, y int64) int64 {
	p := int64(1)
	for y > 0 {
		if y&1 != 0 {
			p *= x
		}
		y >>= 1
		x *= x
	}
	return p
}

// Pow ...
func Pow(args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}

	if x, ok := args[0].(*Integer); ok {
		if y, ok := args[1].(*Integer); ok {
			value := pow(x.Value, y.Value)
			return &Integer{Value: value}
		} else {
			return newError("expected argument #2 to `divmod` to be `int` got=%s", args[1].Type())
		}
	} else {
		return newError("expected argument #1 to `divmod` to be `int` got=%s", args[0].Type())
	}
}
