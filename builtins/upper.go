package builtins

import (
	"strings"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Upper ...
func Upper(args ...object.Object) object.Object {
	if err := typing.Check(
		"upper", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.STRING),
	); err != nil {
		return newError(err.Error())
	}

	return &object.String{Value: strings.ToUpper(args[0].(*object.String).Value)}
}
