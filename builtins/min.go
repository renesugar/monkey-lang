package builtins

import (
	"sort"

	. "github.com/prologic/monkey-lang/object"
)

// Min ...
func Min(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if a, ok := args[0].(*Array); ok {
		// TODO: Make this more generic
		xs := make([]int, len(a.Elements))
		for n, e := range a.Elements {
			if i, ok := e.(*Integer); ok {
				xs = append(xs, int(i.Value))
			} else {
				return newError("item #%d  not an `int` got=%s", n, e.Type())
			}
		}
		sort.Ints(xs)
		return &Integer{Value: int64(xs[0])}
	}
	return newError("argument #1 to `min` expected to be `array` got=%T", args[0].Type())
}
