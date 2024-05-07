package store

import (
	"fmt"
	"sync"
	"time"
)

type Status int

const (
	Recieved Status = iota
	Processing
	Processed
)

type Message struct {
	Topic     string
	ProcessID string

	Data     []byte
	StoredAt time.Time

	Status Status
}

func NewMessage(topic, processID string, data []byte) *Message {
	return &Message{
		Topic:     topic,
		ProcessID: processID,
		Data:      data,
		StoredAt:  time.Now(),
		Status:    Recieved,
	}
}

func (m *Message) Process() {
	m.Status = Processing
}

func (m *Message) Done() {
	m.Status = Processed
}

func (m *Message) String() string {
	return fmt.Sprintf("Message{Topic: %s, ProcessID: %s, Data: %v, StoredAt: %s}", m.Topic, m.ProcessID, m.Data, m.StoredAt)
}

type InlineMemoryStore struct {
	mux sync.Mutex

	maxDataSizePerMessage int

	stream []*Message
}

type Store interface {
	Push(message *Message)
	Stream() []*Message
}

var _ Store = (*InlineMemoryStore)(nil)

func NewStore(maxDataSizePerMessage int) *InlineMemoryStore {
	return &InlineMemoryStore{
		mux:                   sync.Mutex{},
		maxDataSizePerMessage: maxDataSizePerMessage,
		stream:                make([]*Message, 0),
	}
}

func (s *InlineMemoryStore) Push(message *Message) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.stream = append(s.stream, message)
}

func (s *InlineMemoryStore) Stream() []*Message {
	return s.stream
}

func (s *InlineMemoryStore) String() string {
	ret := ""
	for _, message := range s.stream {
		ret += message.String() + "\n"
	}

	return ret
}
