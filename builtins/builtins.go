package builtins

import (
	"fmt"
	"sort"

	. "github.com/prologic/monkey-lang/object"
)

// Builtins ...
var Builtins = map[string]*Builtin{
	"len":    &Builtin{Name: "len", Fn: Len},
	"input":  &Builtin{Name: "input", Fn: Input},
	"print":  &Builtin{Name: "print", Fn: Print},
	"first":  &Builtin{Name: "first", Fn: First},
	"last":   &Builtin{Name: "last", Fn: Last},
	"rest":   &Builtin{Name: "rest", Fn: Rest},
	"push":   &Builtin{Name: "push", Fn: Push},
	"pop":    &Builtin{Name: "pop", Fn: Pop},
	"exit":   &Builtin{Name: "exit", Fn: Exit},
	"assert": &Builtin{Name: "assert", Fn: Assert},
	"bool":   &Builtin{Name: "bool", Fn: Bool},
	"int":    &Builtin{Name: "int", Fn: Int},
	"str":    &Builtin{Name: "str", Fn: Str},
	"type":   &Builtin{Name: "type", Fn: TypeOf},
	"args":   &Builtin{Name: "args", Fn: Args},
	"lower":  &Builtin{Name: "lower", Fn: Lower},
	"upper":  &Builtin{Name: "upper", Fn: Upper},
	"join":   &Builtin{Name: "join", Fn: Join},
	"split":  &Builtin{Name: "split", Fn: Split},
	"find":   &Builtin{Name: "find", Fn: Find},
	"read":   &Builtin{Name: "read", Fn: Read},
	"write":  &Builtin{Name: "write", Fn: Write},
	"ffi":    &Builtin{Name: "ffi", Fn: FFI},
	"abs":    &Builtin{Name: "abs", Fn: Abs},
	"bin":    &Builtin{Name: "bin", Fn: Bin},
	"hex":    &Builtin{Name: "hex", Fn: Hex},
	"ord":    &Builtin{Name: "ord", Fn: Ord},
	"chr":    &Builtin{Name: "chr", Fn: Chr},
	"divmod": &Builtin{Name: "divmod", Fn: Divmod},
	"hash":   &Builtin{Name: "hash", Fn: HashOf},
	"id":     &Builtin{Name: "id", Fn: IdOf},
	"oct":    &Builtin{Name: "oct", Fn: Oct},
	"pow":    &Builtin{Name: "pow", Fn: Pow},
	"min":    &Builtin{Name: "min", Fn: Min},
	"max":    &Builtin{Name: "max", Fn: Max},
}

// BuiltinsIndex ...
var BuiltinsIndex []*Builtin

func init() {
	var keys []string
	for k := range Builtins {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		BuiltinsIndex = append(BuiltinsIndex, Builtins[k])
	}
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
