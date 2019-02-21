package object

import (
	"bytes"
	"strings"
)

// Array is the array literal type that holds a slice of Object(s)
type Array struct {
	Elements []Object
}

func (a *Array) Copy() *Array {
	elements := make([]Object, len(a.Elements))
	for i, e := range a.Elements {
		elements[i] = e
	}
	return &Array{Elements: elements}
}

func (a *Array) Reverse() {
	for i, j := 0, len(a.Elements)-1; i < j; i, j = i+1, j-1 {
		a.Elements[i], a.Elements[j] = a.Elements[j], a.Elements[i]
	}
}

func (a *Array) Len() int {
	return len(a.Elements)
}

func (a *Array) Swap(i, j int) {
	a.Elements[i], a.Elements[j] = a.Elements[j], a.Elements[i]
}

func (a *Array) Less(i, j int) bool {
	if cmp, ok := a.Elements[i].(Comparable); ok {
		return cmp.Compare(a.Elements[j]) == -1
	}
	return false
}

func (ao *Array) Compare(other Object) int {
	if obj, ok := other.(*Array); ok {
		if len(ao.Elements) != len(obj.Elements) {
			return -1
		}
		for i, el := range ao.Elements {
			cmp, ok := el.(Comparable)
			if !ok {
				return -1
			}
			if cmp.Compare(obj.Elements[i]) != 0 {
				return cmp.Compare(obj.Elements[i])
			}
		}

		return 0
	}
	return -1
}

func (ao *Array) String() string {
	return ao.Inspect()
}

// Type returns the type of the object
func (ao *Array) Type() Type { return ARRAY }

// Inspect returns a stringified version of the object for debugging
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
