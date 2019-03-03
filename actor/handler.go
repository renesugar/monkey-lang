package actor

type Handler interface {
	Execute(msg Message) error
}
