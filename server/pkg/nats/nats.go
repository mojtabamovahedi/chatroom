package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

const BaseSubject = "snapp.chatroom"

type Nats struct {
	conn *nats.Conn
}


// New creates a new NATS connection
func New(host string, port int) (*Nats, error) {
    opts := []nats.Option{
        nats.Timeout(5 * time.Second),       // Set connection timeout
        nats.ReconnectWait(2 * time.Second), // Set reconnect wait time
        nats.MaxReconnects(10),              // Set max reconnect attempts
    }

    // Connect to NATS server
    conn, err := nats.Connect(fmt.Sprintf("nats://%s:%d", host, port), opts...)
    if err != nil {
        return nil, err
    }

    return &Nats{conn: conn}, nil
}

// Publish sends a message to a subject
func (n *Nats) Publish(subject string, data []byte) error {
    if n.conn == nil {
        return nats.ErrInvalidConnection
    }
    return n.conn.Publish(subject, data)
}

// Subscribe subscribes to a subject with a message handler
func (n *Nats) Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
    if n.conn == nil {
        return nil, nats.ErrInvalidConnection
    }
    return n.conn.Subscribe(subject, handler)
}

// Request sends a request and waits for a response
func (n *Nats) Request(subject string, data []byte, timeout time.Duration) (*nats.Msg, error) {
    if n.conn == nil {
        return nil, nats.ErrInvalidConnection
    }
    return n.conn.Request(subject, data, timeout)
}

// Drain gracefully closes the connection allowing pending messages to be sent
func (n *Nats) Drain() error {
    if n.conn == nil {
        return nats.ErrInvalidConnection
    }
    return n.conn.Drain()
}

// Close closes the NATS connection
func (n *Nats) Close() {
    if n.conn != nil {
        n.conn.Close()
    }
}
