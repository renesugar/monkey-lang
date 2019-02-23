package builtins

import (
	"fmt"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Print ...
func Print(args ...object.Object) object.Object {
	if err := typing.Check(
		"print", args,
		typing.MinimumArgs(1),
		typing.WithTypes(object.STRING),
	); err != nil {
		return newError(err.Error())
	}

	var s string
	if len(args) == 1 {
		s = args[0].(*object.String).Value
	}

	fmt.Println(s)

	return nil
}
