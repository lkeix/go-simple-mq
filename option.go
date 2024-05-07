package gosimplemq

type ClientOption func(*MQClient)

func WithMaxMessageSize(size int) ClientOption {
	return func(c *MQClient) {
		c.maxMessageSize = size
	}
}
