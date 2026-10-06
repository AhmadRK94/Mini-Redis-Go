package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

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

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if scanner.Err() != nil {
			conn.Write([]byte(scanner.Err().Error()))
		}
		line := scanner.Text()
		inputs := strings.Fields(line)
		cmd := inputs[0]
		args := inputs[1:]
		switch cmd {
		case "GET":
			if len(args) > 1 {
				conn.Write([]byte("invalid input\n"))
			}
			v, err := s.store.Get(args[0])
			if err != nil {
				conn.Write([]byte(err.Error()))
			}
			conn.Write([]byte(fmt.Sprintf("%s\n", v)))
		case "SETNX":
			if len(args) > 2 {
				conn.Write([]byte("invalid input\n"))
			}
			err := s.store.SetNX(args[0], args[1])
			if err != nil {
				conn.Write([]byte(err.Error()))
			}
			conn.Write([]byte("Data stored Succesfully.\n"))
		case "SET":
			if len(args) > 2 {
				conn.Write([]byte("invalid input"))
			}
			err := s.store.Set(args[0], args[1])
			if err != nil {
				conn.Write([]byte(err.Error()))
			}
			conn.Write([]byte("Data stored Succesfully.\n"))
		}
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
