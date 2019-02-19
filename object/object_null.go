package object

// Null is the null type and used to represent the absence of a value
type Null struct{}

func (n *Null) Equal(other Object) bool {
	_, ok := other.(*Null)
	return ok
}

func (n *Null) String() string {
	return n.Inspect()
}

// Type returns the type of the object
func (n *Null) Type() Type { return NULL }

// Inspect returns a stringified version of the object for debugging
func (n *Null) Inspect() string { return "null" }
