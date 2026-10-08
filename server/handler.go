package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic in connection handler: %v\n", r)
		}
	}()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		inputs := strings.Fields(line)
		if len(inputs) == 0 {
			continue
		}
		cmd := strings.ToUpper(inputs[0])
		args := inputs[1:]
		switch cmd {
		case "GET":
			if len(args) != 1 {
				fmt.Fprintln(conn, "Error invalid number of arguments")
				continue
			}
			v, err := s.store.Get(args[0])
			if err != nil {
				fmt.Fprintf(conn, "ERROR %s\n", err)
				continue
			}
			fmt.Fprintf(conn, "%s\n", v)
		case "SETNX":
			if len(args) != 2 {
				fmt.Fprintln(conn, "Error invalid number of arguments")
				continue
			}
			err := s.store.SetNX(args[0], args[1])
			if err != nil {
				fmt.Fprintf(conn, "ERROR %s\n", err)
				continue
			}
			fmt.Fprintln(conn, "OK")
		case "SET":
			if len(args) != 2 {
				fmt.Fprintln(conn, "Error invalid number of arguments")
				continue
			}
			err := s.store.Set(args[0], args[1])
			if err != nil {
				fmt.Fprintf(conn, "ERROR %s\n", err)
				continue
			}
			fmt.Fprintln(conn, "OK")
		case "SETTTL":
			if len(args) != 3 {
				fmt.Fprintln(conn, "Error invalid number of arguments")
				continue
			}
			ttl, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Fprintf(conn, "Error invalid expiration time: %s\n", err.Error())
				continue
			}
			err = s.store.SetWithTTL(args[0], args[1], time.Duration(ttl)*time.Second)
			if err != nil {
				fmt.Fprintln(conn, "Error:", err)
				continue
			}
			fmt.Fprintln(conn, "OK")
		case "DEL":
			if len(args) != 1 {
				fmt.Fprintln(conn, "Error invalid number of arguments")
				continue
			}
			err := s.store.Delete(args[0])
			if err != nil {
				fmt.Fprintf(conn, "ERROR %s\n", err)
				continue
			}
			fmt.Fprintln(conn, "OK")
		case "TTL":
			if len(args) != 1 {
				fmt.Fprintln(conn, "Error invalid number of arguments")
				continue
			}
			ttl := s.store.GetTTL(args[0])
			fmt.Fprintln(conn, ttl)
		case "EXPIRE":
			if len(args) != 2 {
				fmt.Fprintln(conn, "Error invalid number of arguments")
				continue
			}
			ttl, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Fprintf(conn, "Error invalid expiration time: %s\n", err.Error())
				continue
			}
			err = s.store.SetExpire(args[0], time.Duration(ttl)*time.Second)
			if err != nil {
				fmt.Fprintln(conn, "Error:", err)
				continue
			}
			fmt.Fprintln(conn, "OK")
		default:
			fmt.Fprintln(conn, "ERROR unknown command")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("connection error: %v", err)
	}

}
