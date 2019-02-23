package builtins

import (
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// TypeOf ...
func TypeOf(args ...object.Object) object.Object {
	if err := typing.Check(
		"type", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	return &object.String{Value: string(args[0].Type())}
}
