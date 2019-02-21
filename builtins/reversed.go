package builtins

import (
	. "github.com/prologic/monkey-lang/object"
)

// Reversed ...
func Reversed(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if a, ok := args[0].(*Array); ok {
		newArray := a.Copy()
		newArray.Reverse()
		return newArray
	}
	return newError("argument #1 to `reversed` expected to be `array` got=%T", args[0].Type())
}
