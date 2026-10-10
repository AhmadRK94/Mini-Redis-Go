package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"

	"www.github.com/AhmadRK94/nanoredis/store"
)

func TestHandleConnection(t *testing.T) {
	st := store.NewStore()
	server := NewServer(st)

	client, serverConn := net.Pipe()
	done := make(chan struct{})

	go func() {
		defer close(done)
		server.handleConnection(serverConn)
	}()

	reader := bufio.NewReader(client)

	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name:     "SET",
			command:  "SET name ahmad",
			expected: "OK\n",
		},
		{
			name:     "GET existing key",
			command:  "GET name",
			expected: "ahmad\n",
		},
		{
			name:     "TTL without expiration",
			command:  "TTL name",
			expected: "-1\n",
		},
		{
			name:     "unknown command",
			command:  "INVALID",
			expected: "ERROR unknown command\n",
		},
		{
			name:     "invalid argument count",
			command:  "GET",
			expected: "Error invalid number of arguments\n",
		},
		{
			name:     "DEL",
			command:  "DEL name",
			expected: "OK\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := fmt.Fprintln(client, tt.command); err != nil {
				t.Fatalf("failed to send command: %v", err)
			}

			got, err := reader.ReadString('\n')
			if err != nil {
				t.Fatalf("failed to read response: %v", err)
			}

			if got != tt.expected {
				t.Errorf(
					"response = %q, want %q",
					strings.TrimSpace(got),
					strings.TrimSpace(tt.expected),
				)
			}
		})
	}

	_ = client.Close()
	<-done
}
