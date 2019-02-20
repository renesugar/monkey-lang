package builtins

import (
	"fmt"

	. "github.com/prologic/monkey-lang/object"
)

// Print ...
func Print(args ...Object) Object {
	for _, arg := range args {
		fmt.Println(arg.String())
	}

	return nil
}
