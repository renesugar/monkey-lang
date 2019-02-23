package builtins

import (
	"fmt"
	"strconv"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Bin ...
func Bin(args ...object.Object) object.Object {
	if err := typing.Check(
		"bin", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.INTEGER),
	); err != nil {
		return newError(err.Error())
	}

	i := args[0].(*object.Integer)
	return &object.String{Value: fmt.Sprintf("0b%s", strconv.FormatInt(i.Value, 2))}
}
