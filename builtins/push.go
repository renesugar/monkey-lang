package builtins

import (
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Push ...
func Push(args ...object.Object) object.Object {
	if err := typing.Check(
		"push", args,
		typing.ExactArgs(2),
		typing.WithTypes(object.ARRAY),
	); err != nil {
		return newError(err.Error())
	}

	arr := args[0].(*object.Array)
	newArray := arr.Copy()
	newArray.Append(args[1])
	return newArray
}
