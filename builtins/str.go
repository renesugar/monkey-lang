package builtins

import (
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Str ...
func Str(args ...object.Object) object.Object {
	if err := typing.Check(
		"str", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	return &object.String{Value: args[0].String()}
}
