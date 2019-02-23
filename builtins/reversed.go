package builtins

import (
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Reversed ...
func Reversed(args ...object.Object) object.Object {
	if err := typing.Check(
		"reversed", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.ARRAY),
	); err != nil {
		return newError(err.Error())
	}

	arr := args[0].(*object.Array)
	newArray := arr.Copy()
	newArray.Reverse()
	return newArray
}
