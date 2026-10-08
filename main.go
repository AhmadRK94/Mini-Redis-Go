package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"www.github.com/AhmadRK94/Mini-Redis-Go/server"
	"www.github.com/AhmadRK94/Mini-Redis-Go/store"
)

func main() {
	st := store.NewStore()

	if err := st.Load("data_dump.json"); err != nil {
		log.Printf("Could not load data: %v", err)
	}

	stopCleanup := make(chan struct{})
	st.StartCleanup(time.Second, stopCleanup)

	srv := server.NewServer(st, ":6379")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	go func() {
		<-signals
		log.Println("Shutting down...")
		_ = srv.Shutdown()
	}()

	if err := srv.Start(); err != nil {
		log.Printf("Server stopped: %v", err)
	}

	close(stopCleanup)

	if err := st.Save("data_dump.json"); err != nil {
		log.Printf("Could not save data: %v", err)
	} else {
		log.Println("Data saved.")
	}
}
