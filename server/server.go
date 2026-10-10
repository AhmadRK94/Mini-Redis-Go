package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"www.github.com/AhmadRK94/Mini-Redis-Go/store"
)

type Server struct {
	store    *store.Store
	listener net.Listener
}

func NewServer(store *store.Store) *Server {
	return &Server{
		store: store,
	}
}

func (s *Server) start(port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	s.listener = listener
	defer listener.Close()

	fmt.Printf("Server running at localhost:%s...\n", port)

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

func (s *Server) shutdown() error {
	if s.listener == nil {
		return nil
	}
	return s.listener.Close()
}

func (s *Server) ListenAndServe(port string) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- s.start(port)
	}()

	select {
	case sig := <-signals:
		log.Printf("Received %v; shutting down...", sig)

		if err := s.shutdown(); err != nil {
			return err
		}

		return <-serverErr

	case err := <-serverErr:
		return err
	}
}
