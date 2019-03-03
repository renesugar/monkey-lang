package actor

import (
	"fmt"
	"log"
)

const (
	Created State = iota
	Started
	Running
	Stopped
	Crashed

	DefaultMailboxSize = 64 // 64 in-flight messages per Actor mailbox
)

type State int

func (s State) String() string {
	switch s {
	case Created:
		return "CREATED"
	case Started:
		return "STARTED"
	case Running:
		return "RUNNING"
	case Stopped:
		return "STOPPED"
	case Crashed:
		return "CRASHED"
	default:
		return "UNKNOWN"
	}
}

// Actor is the actor type used to represent a single actor (Actor Model) in
// an actor system. It has an internal queue (mailbox) and state (state).
// The actor also holds a user exposed state (State) and address  (Address).
type Actor struct {
	mailbox Mailbox
	state   State
	handler Handler
	address Address
}

// NewActor ...
func NewActor(handler Handler) *Actor {
	return &Actor{
		handler: handler,
		mailbox: NewMailbox(DefaultMailboxSize),
		address: NewAddress(),
	}
}

func (a *Actor) String() string {
	return fmt.Sprintf("<actor %s %s>", a.address, a.state)
}

func (a *Actor) run() error {
	var err error

	defer func() {
		if x := recover(); x != nil {
			log.Printf("run time panic: %v", x)
			a.state = Crashed
			// TODO: Notify parent actor
			a.Send(Message{Flags: Event, Payload: []byte("crashed")})
		} else if err != nil {
			log.Printf("error handling message: %s", err)
			a.state = Crashed
			// TODO: Notify parent actor
			a.Send(Message{Flags: Event, Payload: []byte("crashed")})
		} else {
			log.Printf("%s stopped", a)
			a.state = Stopped
			// TODO: Notify parent actor
			a.Send(Message{Flags: Event, Payload: []byte("stopped")})
		}
	}()

	log.Printf("%s stared", a)

	a.state = Started
	for {
		msg := a.mailbox.Get()
		log.Printf("%s received %#v", a, msg)
		if (msg.Flags & Stop) != 0 {
			return nil
		}

		err = a.handler.Execute(msg)
		if err != nil {
			return err
		}
	}
}

func (a *Actor) Send(msg Message) {
	a.mailbox.Put(msg)
}

func (a *Actor) Start() {
	defer func() {
		// TODO: Notify parent actor
		a.Send(Message{Flags: Event, Payload: []byte("started")})
	}()
	go a.run()
}

func (a *Actor) Stop() {
	a.Send(Message{Flags: Stop})
}
