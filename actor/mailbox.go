package actor

type Mailbox interface {
	Get() Message
	Put(msg Message)
}

var NewMailbox = NewChannelMailbox

type ChannelMailbox struct {
	ch chan Message
}

func NewChannelMailbox(n int) *ChannelMailbox {
	return &ChannelMailbox{ch: make(chan Message, n)}
}

func (m *ChannelMailbox) Get() Message {
	msg := <-m.ch
	return msg
}

func (m *ChannelMailbox) Put(msg Message) {
	m.ch <- msg
}
