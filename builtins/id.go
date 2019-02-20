package builtins

import (
	"fmt"

	. "github.com/prologic/monkey-lang/object"
)

// IdOf ...
func IdOf(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	arg := args[0]

	if n, ok := arg.(*Null); ok {
		return &String{Value: fmt.Sprintf("%p", n)}
	} else if b, ok := arg.(*Boolean); ok {
		return &String{Value: fmt.Sprintf("%p", b)}
	} else if i, ok := arg.(*Integer); ok {
		return &String{Value: fmt.Sprintf("%p", i)}
	} else if s, ok := arg.(*String); ok {
		return &String{Value: fmt.Sprintf("%p", s)}
	} else if a, ok := arg.(*Array); ok {
		return &String{Value: fmt.Sprintf("%p", a)}
	} else if h, ok := arg.(*Hash); ok {
		return &String{Value: fmt.Sprintf("%p", h)}
	} else if f, ok := arg.(*Function); ok {
		return &String{Value: fmt.Sprintf("%p", f)}
	} else if c, ok := arg.(*Closure); ok {
		return &String{Value: fmt.Sprintf("%p", c)}
	} else if b, ok := arg.(*Builtin); ok {
		return &String{Value: fmt.Sprintf("%p", b)}
	} else {
		return newError("argument 31 to `id()` unsupported got=%T", arg.Type())
	}
}
