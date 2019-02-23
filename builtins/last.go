package builtins

import (
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Last ...
func Last(args ...object.Object) object.Object {
	if err := typing.Check(
		"last", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.ARRAY),
	); err != nil {
		return newError(err.Error())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}

	return nil
}
