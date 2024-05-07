package gosimplemq

import (
	"encoding/json"
	"net"

	"github.com/lkeix/go-simple-mq/internal"
)

type Topic string

type MQClient struct {
	conn net.Conn

	maxMessageSize int
}

type Client interface {
	Send(message []byte) error
	Receive() ([]byte, error)

	ListTopics() ([]Topic, error)
	Publish(topics []Topic, message []byte) error
	Subscribe(topics []Topic, callback func(message []byte) error) error
}

var _ Client = (*MQClient)(nil)

func NewClient(network, addr string, opts ...ClientOption) (*MQClient, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	c := &MQClient{
		conn:           conn,
		maxMessageSize: 1024,
	}
	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *MQClient) Send(message []byte) error {
	m := internal.NewMessage("", message)
	b, _ := json.Marshal(m)

	_, err := c.conn.Write(b)
	return err
}

func (c *MQClient) Receive() ([]byte, error) {
	buf := make([]byte, c.maxMessageSize)
	n, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

func (c *MQClient) ListTopics() ([]Topic, error) {
	return nil, nil
}

func (c *MQClient) Publish(topics []Topic, message []byte) error {
	return nil
}

func (c *MQClient) Subscribe(topics []Topic, callback func(message []byte) error) error {
	return nil
}
