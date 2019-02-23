package builtins

import (
	"io/ioutil"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// Read ...
func Read(args ...object.Object) object.Object {
	if err := typing.Check(
		"read", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.STRING),
	); err != nil {
		return newError(err.Error())
	}

	filename := args[0].(*object.String).Value
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return newError("IOError: error reading from file %s: %s", filename, err)
	}

	return &object.String{Value: string(data)}
}
