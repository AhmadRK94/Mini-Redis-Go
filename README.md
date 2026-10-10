# Nano Redis — A Redis-Inspired Key-Value Store in Go

A lightweight, Redis-inspired in-memory key-value store built from scratch in Go. The project explores TCP networking, concurrent programming, key expiration, and data persistence while implementing a simple text-based command protocol.

## Features

* **TCP server:** Accepts multiple client connections using Go's `net` package.
* **In-memory storage:** Stores key-value pairs in a shared map.
* **Concurrent access:** Uses synchronization primitives to protect shared state.
* **Key expiration (TTL):** Supports expiration of keys after a specified duration.
* **Background cleanup:** Periodically removes expired entries.
* **Data persistence:** Loads and saves stored data using JSON snapshots.
* **Command-based protocol:** Allows clients to interact with the store through text commands.

## Tech Stack

* **Language:** Go
* **Networking:** `net`
* **Concurrency:** Goroutines, channels, `sync`
* **Data storage:** Go maps
* **Persistence:** JSON

## Getting Started

### Prerequisites

* Go installed on your machine.
* A terminal or TCP client to connect to the server.

### Installation

Clone the repository:

```bash
git clone https://github.com/AhmadRK94/nanoredis.git
cd nanoredis
```

Initialize dependencies if needed:

```bash
go mod tidy
```

### Run the server

```bash
go run .
```

The server listens on port `6379` by default.

To build an executable:

```bash
go build -o nanoredis.exe .
```

Then run the generated executable.

## Commands

The server uses a simple, text-based command protocol. Commands are sent as lines of text, with arguments separated by spaces.

| Command  | Description                         | Example            |
| -------- | ----------------------------------- | ------------------ |
| `SET`    | Store a key-value pair              | `SET name ahmad`   |
| `SETNX`  | Set a key only if it doesn't exist  | `SETNX name ahmad` |
| `SETTTL` | Store a key-value pair with TTL     | `SETTTL name ahmad 60`|
| `GET`    | Retrieve a value                    | `GET name`         |
| `DEL`    | Delete a key                        | `DEL name`         |
| `EXPIRE` | Set a key's expiration in seconds   | `EXPIRE name 60`   |
| `TTL`    | Get the remaining lifetime of a key | `TTL name`         |

### Connect to the server

You can use `netcat` if it is installed:

```bash
nc localhost 6379
```

Or use `telnet`:

```bash
telnet localhost 6379
```

Then enter commands, for example:

```text
SET name ahmad
GET name
TTL name
```

The exact responses depend on the protocol implemented by the server.

## Persistence

The project uses a JSON snapshot file named `data_dump.json` to persist stored data between server runs.

The intended lifecycle is:

1. Load previously saved data during startup.
2. Serve client requests and update the in-memory store.
3. Save the store during graceful shutdown.

To preserve data reliably, the server must complete its shutdown and save operations. Forcefully terminating the process may prevent the final snapshot from being written.

## Project Structure

The exact structure may evolve as the project grows. A possible layout is:

* **`main.go`** — Initializes the store, starts background cleanup, runs the server, and coordinates shutdown.
* **`server/`** — Handles TCP connections, parses commands, and sends responses.
* **`store/`** — Implements key-value operations, synchronization, expiration logic, and save and loading the snapshot.
* **`*_test.go`** — Contains automated tests.

## Testing

Run all tests:

```bash
go test ./...
```

Show detailed test output:

```bash
go test -v ./...
```

Check for data races:

```bash
go test -race ./...
```

The race detector is particularly useful for testing concurrent access to the shared store.

## Design Goals

This project is intended to explore several core backend engineering concepts:

* TCP client-server communication.
* Concurrent connection handling with goroutines.
* Safe shared-memory access using mutexes.
* TTL management and background workers.
* Persistence and graceful shutdown.
* Unit testing and concurrent system validation.

## Current Limitations

This is an educational project inspired by Redis, not a production-ready replacement. Depending on the implementation, it may not include Redis protocol compatibility, authentication, replication, advanced data structures, or crash-safe persistence.
