package builtins

import (
	"syscall"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// DefaultBufferSize is the default buffer size
const DefaultBufferSize = 4096

// Read ...
func Read(args ...object.Object) object.Object {
	if err := typing.Check(
		"read", args,
		typing.RangeOfArgs(1, 2),
		typing.WithTypes(object.INTEGER, object.INTEGER),
	); err != nil {
		return newError(err.Error())
	}

	var (
		fd int
		n  = DefaultBufferSize
	)

	fd = int(args[0].(*object.Integer).Value)

	if len(args) == 2 {
		n = int(args[1].(*object.Integer).Value)
	}

	buf := make([]byte, n)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return newError("IOError: %s", err)
	}

	return &object.String{Value: string(buf[:n])}
}
