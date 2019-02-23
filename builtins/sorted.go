package builtins

import (
	"sort"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Sorted ...
func Sorted(args ...object.Object) object.Object {
	if err := typing.Check(
		"sort", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.ARRAY),
	); err != nil {
		return newError(err.Error())
	}

	arr := args[0].(*object.Array)
	newArray := arr.Copy()
	sort.Sort(newArray)
	return newArray
}
