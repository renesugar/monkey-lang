package builtins

import (
	"syscall"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Accept ...
func Accept(args ...object.Object) object.Object {
	if err := typing.Check(
		"accept", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.INTEGER),
	); err != nil {
		return newError(err.Error())
	}

	var (
		nfd int
		err error
	)

	fd := int(args[0].(*object.Integer).Value)

	nfd, _, err = syscall.Accept(fd)
	if err != nil {
		return newError("SocketError: %s", err)
	}

	return &object.Integer{Value: int64(nfd)}
}
