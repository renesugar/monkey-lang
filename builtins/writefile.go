package builtins

import (
	"io/ioutil"

	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/typing"
)

// WriteFile ...
func WriteFile(args ...object.Object) object.Object {
	if err := typing.Check(
		"writefile", args,
		typing.ExactArgs(2),
		typing.WithTypes(object.STRING, object.STRING),
	); err != nil {
		return newError(err.Error())
	}

	filename := args[0].(*object.String).Value
	data := []byte(args[1].(*object.String).Value)

	err := ioutil.WriteFile(filename, data, 0755)
	if err != nil {
		return newError("IOError: error writing file %s: %s", filename, err)
	}

	return &object.Null{}
}
