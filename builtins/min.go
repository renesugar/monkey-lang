package builtins

import (
	"sort"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Min ...
func Min(args ...object.Object) object.Object {
	if err := typing.Check(
		"min", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.ARRAY),
	); err != nil {
		return newError(err.Error())
	}

	a := args[0].(*object.Array)
	// TODO: Make this more generic
	xs := make([]int, len(a.Elements))
	for n, e := range a.Elements {
		if i, ok := e.(*object.Integer); ok {
			xs = append(xs, int(i.Value))
		} else {
			return newError("item #%d  not an `int` got=%s", n, e.Type())
		}
	}
	sort.Ints(xs)
	return &object.Integer{Value: int64(xs[0])}
}
