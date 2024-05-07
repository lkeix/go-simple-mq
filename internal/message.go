package internal

import "github.com/google/uuid"

type Message struct {
	Topic     string `json:"topic"`
	ProcessID string `json:"process_id"`

	Data []byte `json:"data"`
}

func NewMessage(topic string, data []byte) *Message {
	processID := uuid.New().String()

	return &Message{
		Topic:     topic,
		ProcessID: processID,
		Data:      data,
	}
}
