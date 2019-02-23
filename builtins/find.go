package builtins

import (
	"sort"
	"strings"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Find ...
func Find(args ...object.Object) object.Object {
	if err := typing.Check(
		"find", args,
		typing.ExactArgs(2),
	); err != nil {
		return newError(err.Error())
	}

	// find("foobar", "bo")
	if haystack, ok := args[0].(*object.String); ok {
		if err := typing.Check(
			"find", args,
			typing.WithTypes(object.STRING, object.STRING),
		); err != nil {
			return newError(err.Error())
		}

		needle := args[1].(*object.String)
		index := strings.Index(haystack.Value, needle.Value)
		return &object.Integer{Value: int64(index)}
	}

	// find([1, 2, 3], 2)
	if haystack, ok := args[0].(*object.Array); ok {
		needle := args[1].(object.Comparable)
		i := sort.Search(len(haystack.Elements), func(i int) bool {
			return needle.Compare(haystack.Elements[i]) == 0
		})
		if i < len(haystack.Elements) && needle.Compare(haystack.Elements[i]) == 0 {
			return &object.Integer{Value: int64(i)}
		}
		return &object.Integer{Value: -1}
	}

	return newError(
		"TypeError: find() expected argument #1 to be `array` or `str` got `%s`",
		args[0].Type(),
	)
}
