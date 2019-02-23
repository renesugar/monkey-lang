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
	); err != nil {
		return newError(err.Error())
	}

	fmt.Println(args[0].String())

	return nil
}
