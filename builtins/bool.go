package builtins

import (
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Bool ...
func Bool(args ...object.Object) object.Object {
	if err := typing.Check(
		"bool", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	return &object.Boolean{Value: args[0].Bool()}
}
