package main

import (
	"log"
	"time"

	"www.github.com/AhmadRK94/Mini-Redis-Go/server"
	"www.github.com/AhmadRK94/Mini-Redis-Go/store"
)

func main() {
	store := store.NewStore()
	stopCleanup := make(chan struct{})
	store.StartCleanup(time.Second, stopCleanup)
	server := server.NewServer(store, ":6379")
	log.Fatal(server.Start())
	close(stopCleanup)
}
