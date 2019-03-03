package actor

type Flag int

const (
	Stop Flag = iota
	Event
)

type Payload []byte

func (p Payload) String() string {
	return string(p)
}

type Message struct {
	Flags   Flag
	Address Address
	Payload Payload
}

func NewMessage(address Address, payload []byte) Message {
	return Message{Address: address, Payload: payload}
}
