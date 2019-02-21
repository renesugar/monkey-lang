package builtins

import (
	"sort"

	. "github.com/prologic/monkey-lang/object"
)

// Sorted ...
func Sorted(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if a, ok := args[0].(*Array); ok {
		newArray := a.Copy()
		sort.Sort(newArray)
		return newArray
	}
	return newError("argument #1 to `sorted` expected to be `array` got=%T", args[0].Type())
}
