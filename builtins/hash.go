package builtins

import (
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// HashOf ...
func HashOf(args ...object.Object) object.Object {
	if err := typing.Check(
		"hash", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	if hash, ok := args[0].(object.Hashable); ok {
		return &object.Integer{Value: int64(hash.HashKey().Value)}
	}

	return newError("TypeError: hash() expected argument #1 to be hashable")
}
