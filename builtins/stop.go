package builtins

import (
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Stop ...
func Stop(args ...object.Object) object.Object {
	if err := typing.Check(
		"stop", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.ACTOR),
	); err != nil {
		return newError(err.Error())
	}

	actor := args[0].(*object.Actor)
	actor.Ref.Stop()

	return &object.Null{}
}
