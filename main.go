package main

import (
	"log"

	"www.github.com/AhmadRK94/Mini-Redis-Go/server"
	"www.github.com/AhmadRK94/Mini-Redis-Go/store"
)

func main() {
	store := store.NewStore()
	server := server.NewServer(store, ":6379")
	log.Fatal(server.Start())
}
