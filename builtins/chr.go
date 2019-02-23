package builtins

import (
	"fmt"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Chr ...
func Chr(args ...object.Object) object.Object {
	if err := typing.Check(
		"chr", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.INTEGER),
	); err != nil {
		return newError(err.Error())
	}

	i := args[0].(*object.Integer)
	return &object.String{Value: fmt.Sprintf("%c", rune(i.Value))}
}
