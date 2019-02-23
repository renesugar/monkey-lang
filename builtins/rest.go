package builtins

import (
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Rest ...
func Rest(args ...object.Object) object.Object {
	if err := typing.Check(
		"rest", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.ARRAY),
	); err != nil {
		return newError(err.Error())
	}

	arr := args[0].(*object.Array)
	newArray := arr.Copy()
	newArray.PopLeft()
	return newArray
}
