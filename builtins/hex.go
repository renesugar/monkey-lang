package builtins

import (
	"fmt"
	"strconv"

	. "github.com/prologic/monkey-lang/object"
)

// Hex ...
func Hex(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if i, ok := args[0].(*Integer); ok {
		return &String{Value: fmt.Sprintf("0x%s", strconv.FormatInt(i.Value, 16))}
	}
	return newError("argument to `hex` not supported, got %s", args[0].Type())
}
