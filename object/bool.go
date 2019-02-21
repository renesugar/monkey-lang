package object

import (
	"fmt"
)

// Boolean is the boolean type and used to represent boolean literals and
// holds an interval bool value
type Boolean struct {
	Value bool
}

func (b *Boolean) Int() int64 {
	if b.Value {
		return 1
	}
	return 0
}

func (b *Boolean) Less(other Object) bool {
	if obj, ok := other.(*Boolean); ok {
		return b.Int() < obj.Int()
	}
	return false
}

func (b *Boolean) Equal(other Object) bool {
	if obj, ok := other.(*Boolean); ok {
		return b.Value == obj.Value
	}
	return false
}

func (b *Boolean) String() string {
	return b.Inspect()
}

// Clone creates a new copy
func (b *Boolean) Clone() Object {
	return &Boolean{Value: b.Value}
}

// Type returns the type of the object
func (b *Boolean) Type() Type { return BOOLEAN }

// Inspect returns a stringified version of the object for debugging
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
