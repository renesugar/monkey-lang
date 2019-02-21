package builtins

import (
	"sort"
	"strings"

	. "github.com/prologic/monkey-lang/object"
)

// Find ...
func Find(args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}

	if haystack, ok := args[0].(*String); ok {
		if needle, ok := args[1].(*String); ok {
			index := strings.Index(haystack.Value, needle.Value)
			return &Integer{Value: int64(index)}
		} else {
			return newError("expected arg #2 to be `str` got got=%T", args[1])
		}
	} else if haystack, ok := args[0].(*Array); ok {
		needle := args[1].(Comparable)
		i := sort.Search(len(haystack.Elements), func(i int) bool {
			return needle.Compare(haystack.Elements[i]) == 0
		})
		if i < len(haystack.Elements) && needle.Compare(haystack.Elements[i]) == 0 {
			return &Integer{Value: int64(i)}
		}
		return &Integer{Value: -1}
	} else {
		return newError("expected arg #1 to be `str` or `array` got got=%T", args[0])
	}
}
