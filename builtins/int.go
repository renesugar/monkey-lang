package builtins

import (
	"strconv"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Int ...
func Int(args ...object.Object) object.Object {
	if err := typing.Check(
		"int", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	switch arg := args[0].(type) {
	case *object.Boolean:
		if arg.Value {
			return &object.Integer{Value: 1}
		}
		return &object.Integer{Value: 0}
	case *object.Integer:
		return arg
	case *object.String:
		n, err := strconv.ParseInt(arg.Value, 10, 64)
		if err != nil {
			return newError("could not parse string to int: %s", err)
		}
		return &object.Integer{Value: n}
	default:
		return &object.Integer{}
	}
}
