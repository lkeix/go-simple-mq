package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/lkeix/go-simple-mq/internal"
	"github.com/lkeix/go-simple-mq/internal/store"
)

type Server struct {
	store store.Store

	listener net.Listener
	conns    map[net.Addr]net.Conn

	network string
	address string
}

func NewServer(network, address string, store store.Store) (*Server, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}

	return &Server{
		store:    store,
		listener: listener,
		conns:    make(map[net.Addr]net.Conn),
		network:  network,
		address:  address,
	}, nil
}

func (s *Server) Serve(ch chan error) {
	fmt.Printf("Server is listening on %s:%s\n", s.network, s.address)
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)
			ch <- err
		}
	}()

	go s.Publish()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		s.conns[conn.RemoteAddr()] = conn
		go s.handleConnection(conn)
	}
}

func (s *Server) Publish() {
	for {
		messages := s.store.Stream()
		for _, message := range messages {
			if message.Status == store.Recieved {
				message.Process()
				for _, conn := range s.conns {
					b, _ := json.Marshal(message)
					conn.Write(b)
				}
				message.Done()
			}
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		if _, err := conn.Read(buf); err != nil {
			delete(s.conns, conn.RemoteAddr())
			conn.Close()
			return
		}
		buf = bytes.Trim(buf, "\x00")

		var msg internal.Message
		if err := json.Unmarshal(buf, &msg); err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}

		qm := store.NewMessage(msg.Topic, msg.ProcessID, msg.Data)
		s.store.Push(qm)
	}
}

func Run(network, address string, maxDataSizePerMessage int) {
	srv, err := NewServer(network, address, store.NewStore(maxDataSizePerMessage))
	if err != nil {
		fmt.Println("Error creating server:", err)
		return
	}

	defer srv.listener.Close()

	ch := make(chan error, 1)
	go srv.Serve(ch)

	if err := <-ch; err != nil {
		log.Fatal(err)
	}
}
