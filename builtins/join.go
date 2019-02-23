package builtins

import (
	"strings"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Join ...
func Join(args ...object.Object) object.Object {
	if err := typing.Check(
		"join", args,
		typing.ExactArgs(2),
		typing.WithTypes(object.ARRAY, object.STRING),
	); err != nil {
		return newError(err.Error())
	}

	arr := args[0].(*object.Array)
	sep := args[1].(*object.String)
	a := make([]string, len(arr.Elements))
	for i, el := range arr.Elements {
		a[i] = el.String()
	}
	return &object.String{Value: strings.Join(a, sep.Value)}
}
