package object

import (
	"bytes"
	"encoding/gob"

	"github.com/prologic/monkey-lang/actor"
)

// Actor is the actor type used to represent a reference to an Actor
// (Actor Model) in actor.Actor
type Actor struct {
	Ref *actor.Actor
}

func (a *Actor) Send(msg Object) error {
	var buffer bytes.Buffer

	// TODO: How do we get the Address (from) of the sender?
	address := actor.Address{} // sender address

	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(&msg)
	if err != nil {
		return err
	}
	payload := buffer.Bytes()

	a.Ref.Send(actor.NewMessage(address, payload))

	return nil
}

func (a *Actor) Bool() bool {
	return true
}

func (a *Actor) Compare(other Object) int {
	return 1
}

func (a *Actor) String() string {
	return a.Inspect()
}

// Type returns the type of the object
func (a *Actor) Type() Type { return ACTOR }

// Inspect returns a stringified version of the object for debugging
func (a *Actor) Inspect() string { return a.Ref.String() }
