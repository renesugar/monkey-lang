package builtins

import (
	"syscall"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Close ...
func Close(args ...object.Object) object.Object {
	if err := typing.Check(
		"close", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.INTEGER),
	); err != nil {
		return newError(err.Error())
	}

	fd := int(args[0].(*object.Integer).Value)

	err := syscall.Close(fd)
	if err != nil {
		return newError("IOError: %s", err)
	}

	return &object.Null{}
}
