package builtins

import (
	"syscall"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Listen ...
func Listen(args ...object.Object) object.Object {
	if err := typing.Check(
		"listen", args,
		typing.ExactArgs(2),
		typing.WithTypes(object.INTEGER, object.INTEGER),
	); err != nil {
		return newError(err.Error())
	}

	fd := int(args[0].(*object.Integer).Value)
	backlog := int(args[1].(*object.Integer).Value)

	if err := syscall.Listen(fd, backlog); err != nil {
		return newError("SocketError: %s", err)
	}

	return &object.Null{}
}
