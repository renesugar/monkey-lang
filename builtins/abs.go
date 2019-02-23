package builtins

import (
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Abs ...
func Abs(args ...object.Object) object.Object {
	if err := typing.Check(
		"abs", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.INTEGER),
	); err != nil {
		return newError(err.Error())
	}

	i := args[0].(*object.Integer)
	value := i.Value
	if value < 0 {
		value = value * -1
	}
	return &object.Integer{Value: value}
}
