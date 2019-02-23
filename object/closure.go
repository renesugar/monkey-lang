package object

import (
	"fmt"

	"github.com/prologic/monkey-lang/code"
)

// CompiledFunction is the compiled function type that holds the function's
// compiled body as bytecode instructions
type CompiledFunction struct {
	Instructions  code.Instructions
	NumLocals     int
	NumParameters int
}

func (cf *CompiledFunction) Bool() bool {
	return true
}

func (cf *CompiledFunction) String() string {
	return cf.Inspect()
}

// Type returns the type of the object
func (cf *CompiledFunction) Type() Type { return COMPILED_FUNCTION }

// Inspect returns a stringified version of the object for debugging
func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}

// Closure is the closure object type that holds a reference to a compiled
// functions and its free variables
type Closure struct {
	Fn   *CompiledFunction
	Free []Object
}

func (c *Closure) Bool() bool {
	return true
}

func (c *Closure) String() string {
	return c.Inspect()
}

// Type returns the type of the object
func (c *Closure) Type() Type { return CLOSURE }

// Inspect returns a stringified version of the object for debugging
func (c *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", c)
}
