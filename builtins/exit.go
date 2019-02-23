package builtins

import (
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Exit ...
func Exit(args ...object.Object) object.Object {
	if err := typing.Check(
		"exit", args,
		typing.RangeOfArgs(0, 1),
		typing.WithTypes(object.INTEGER),
	); err != nil {
		return newError(err.Error())
	}

	var status int
	if len(args) == 1 {
		status = int(args[0].(*object.Integer).Value)
	}

	object.ExitFunction(status)

	return nil
}
