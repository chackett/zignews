package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

// Nat defines functionality to access the Nats message queue
type Nat struct {
	connection *nats.Conn
}

// New returns a new instance of Nats
func New(conn string) (Nat, error) {
	nc, err := nats.Connect(conn)
	if err != nil {
		return Nat{}, errors.Wrap(err, "nats connect")
	}
	return Nat{
		connection: nc,
	}, nil
}
