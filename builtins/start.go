package builtins

import (
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Start ...
func Start(args ...object.Object) object.Object {
	if err := typing.Check(
		"start", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.ACTOR),
	); err != nil {
		return newError(err.Error())
	}

	actor := args[0].(*object.Actor)
	actor.Ref.Start()

	return &object.Null{}
}
