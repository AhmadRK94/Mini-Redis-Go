package server

import (
	"errors"
	"fmt"
	"net"

	"www.github.com/AhmadRK94/Mini-Redis-Go/store"
)

type Server struct {
	store    *store.Store
	port     string
	listener net.Listener
}

func NewServer(store *store.Store, port string) *Server {
	return &Server{
		store: store,
		port:  port,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	s.listener = listener
	defer listener.Close()

	fmt.Printf("Server running at localhost:%s...\n", s.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) Shutdown() error {
	if s.listener == nil {
		return nil
	}
	return s.listener.Close()
}
