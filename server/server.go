package server

import (
	"fmt"
	"net"

	"www.github.com/AhmadRK94/Mini-Redis-Go/store"
)

type Server struct {
	store *store.Store
	port  string
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
	defer listener.Close()
	fmt.Printf("Server start running at localhost:%s...\n", s.port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(conn)
	}

}
