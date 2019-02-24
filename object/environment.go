package object

// NewEnvironment constructs a new Environment object to hold bindings
// of identifiers to their names
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Environment is an object that holds a mapping of names to bound objets
type Environment struct {
	store  map[string]Object
	parent *Environment
}

// Hash returns a new Hash with the names and values of every value in the
// environment. This is used by the module import system to wrap up the
// evaulated module into an object.
func (e *Environment) Hash() *Hash {
	pairs := make(map[HashKey]HashPair)
	for k, v := range e.store {
		s := &String{Value: k}
		pairs[s.HashKey()] = HashPair{Key: s, Value: v}
	}
	return &Hash{Pairs: pairs}
}

// Clone returns a new Environment with the parent set to the current
// environment (enclosing environment)
func (e *Environment) Clone() *Environment {
	env := NewEnvironment()
	env.parent = e
	return env
}

// Get returns the object bound by name
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.parent != nil {
		obj, ok = e.parent.Get(name)
	}
	return obj, ok
}

// Set stores the object with the given name
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
