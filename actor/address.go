package actor

import (
	"fmt"

	"github.com/lithammer/shortuuid"
)

const (
	Local AddressType = iota
	Remote
)

var NewAddress = NewLocalAddress

type AddressType int

func (at AddressType) String() string {
	switch at {
	case Local:
		return "local"
	case Remote:
		return "remote"
	default:
		return "unknown"
	}
}

type Address struct {
	Type AddressType
	Host string
	Port int
	UUID string
}

func (a Address) String() string {
	if a.Type == Local {
		return fmt.Sprintf("%s://%s", a.Type, a.UUID)
	}
	return fmt.Sprintf("%s://%s:%d/%s", a.Type, a.Host, a.Port, a.UUID)
}

func NewLocalAddress() Address {
	return Address{Type: Local, UUID: shortuuid.New()}
}

func NewRemoteAddress(host string, port int) Address {
	return Address{Type: Remote, Host: host, Port: port, UUID: shortuuid.New()}
}
