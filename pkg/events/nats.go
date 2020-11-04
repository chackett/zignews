package events

// // Nat defines functionality to access the Nats message queue
// type Nat struct {
// 	connection *nats.Conn
// }

// // MessageBus defines functionality to publish messages on a message bus
// type MessageBus interface {
// 	Publish(subject string, msg []byte) error
// }

// // New returns a new instance of Nats
// func New(conn string) (*Nat, error) {
// 	nc, err := nats.Connect(conn)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "nats connect")
// 	}
// 	return &Nat{
// 		connection: nc,
// 	}, nil
// }

// // Publish submits a message to the queue using provided subject
// func (n *Nat) Publish(subject string, msg []byte) error {
// 	return errors.Wrap(n.connection.Publish(subject, msg), "nats publish")
// }
