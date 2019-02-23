package builtins

import (
	"strings"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Lower ...
func Lower(args ...object.Object) object.Object {
	if err := typing.Check(
		"lower", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.STRING),
	); err != nil {
		return newError(err.Error())
	}

	str := args[0].(*object.String)
	return &object.String{Value: strings.ToLower(str.Value)}
}
