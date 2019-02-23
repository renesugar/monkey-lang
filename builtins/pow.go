package builtins

import (
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

func pow(x, y int64) int64 {
	p := int64(1)
	for y > 0 {
		if y&1 != 0 {
			p *= x
		}
		y >>= 1
		x *= x
	}
	return p
}

// Pow ...
func Pow(args ...object.Object) object.Object {
	if err := typing.Check(
		"pow", args,
		typing.ExactArgs(2),
		typing.WithTypes(object.INTEGER, object.INTEGER),
	); err != nil {
		return newError(err.Error())
	}

	x := args[0].(*object.Integer)
	y := args[1].(*object.Integer)
	value := pow(x.Value, y.Value)
	return &object.Integer{Value: value}
}
