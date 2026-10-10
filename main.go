package main

import (
	"log"
	"time"

	"www.github.com/AhmadRK94/Mini-Redis-Go/server"
	"www.github.com/AhmadRK94/Mini-Redis-Go/store"
)

func main() {
	store := store.NewStore()

	if err := store.Load("data_dump.json"); err != nil {
		log.Printf("Could not load data: %v", err)
	}

	store.StartCleanup(time.Second)
	defer store.StopCleanup()

	server := server.NewServer(store)
	if err := server.ListenAndServe(":6379"); err != nil {
		log.Printf("Server stopped: %v", err)
	}

	if err := store.Save("data_dump.json"); err != nil {
		log.Printf("Could not save data: %v", err)
	} else {
		log.Println("Data saved.")
	}
}
