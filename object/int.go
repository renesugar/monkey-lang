package object

import (
	"fmt"
)

// Integer is the integer type used to represent integer literals and holds
// an internal int64 value
type Integer struct {
	Value int64
}

func (i *Integer) Compare(other Object) int {
	if obj, ok := other.(*Integer); ok {
		switch {
		case i.Value < obj.Value:
			return -1
		case i.Value > obj.Value:
			return 1
		default:
			return 0
		}
	}
	return -1
}

func (i *Integer) String() string {
	return i.Inspect()
}

// Clone creates a new copy
func (i *Integer) Clone() Object {
	return &Integer{Value: i.Value}
}

// Type returns the type of the object
func (i *Integer) Type() Type { return INTEGER }

// Inspect returns a stringified version of the object for debugging
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
