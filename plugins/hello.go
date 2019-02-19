package main

import (
	. "github.com/prologic/monkey-lang/object"
)

// Hello ...
func Hello(args ...Object) Object {
	return &String{Value: "Hello World!"}
}
