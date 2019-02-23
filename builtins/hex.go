package builtins

import (
	"fmt"
	"strconv"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Hex ...
func Hex(args ...object.Object) object.Object {
	if err := typing.Check(
		"hex", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.INTEGER),
	); err != nil {
		return newError(err.Error())
	}

	i := args[0].(*object.Integer)
	return &object.String{Value: fmt.Sprintf("0x%s", strconv.FormatInt(i.Value, 16))}
}
