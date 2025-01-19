package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

type Nats struct {
	conn *nats.Conn
}

func New(host string, port int) (*Nats, error) {
	opts := []nats.Option{
		nats.Timeout(5 * time.Second),
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(10),
	}

	conn, err := nats.Connect(fmt.Sprintf("nats://%s:%d", host, port), opts...)
	if err != nil {
		return nil, err
	}

	return &Nats{conn: conn}, nil
}

func (n *Nats) Publish(subject string, data []byte) error {
	if n.conn == nil {
		return nats.ErrInvalidConnection
	}
	return n.conn.Publish(subject, data)
}

func (n *Nats) Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	if n.conn == nil {
		return nil, nats.ErrInvalidConnection
	}
	return n.conn.Subscribe(subject, handler)
}

func (n *Nats) Request(subject string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	if n.conn == nil {
		return nil, nats.ErrInvalidConnection
	}
	return n.conn.Request(subject, data, timeout)
}

func (n *Nats) Drain() error {
	if n.conn == nil {
		return nats.ErrInvalidConnection
	}
	return n.conn.Drain()
}

func (n *Nats) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}
