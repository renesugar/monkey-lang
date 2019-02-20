package builtins

import (
	"fmt"

	. "github.com/prologic/monkey-lang/object"
)

// Chr ...
func Chr(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if i, ok := args[0].(*Integer); ok {
		return &String{Value: fmt.Sprintf("%c", rune(i.Value))}
	}
	return newError("argument to `chr` not supported, got %s", args[0].Type())
}
