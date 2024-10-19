# Go Cache Server

A simple TCP server in Go for managing a key-value cache. This project includes:

- A TCP server that allows clients to interact with the cache over a telnet connection.
- A cache with configurable capacity, TTL, and snapshot persistence using `.gob`.
- Graceful shutdown for managing connections and ensuring data integrity.

## Features

- **Key-Value Cache**: Supports `GET`, `SET`, and `DELETE` commands.
- **Persistence**: Cache is saved to disk as a snapshot using Go's `gob` encoding.
- **Eviction**: Items automatically expire based on TTL and are periodically evicted.
- **Graceful Shutdown**: Handles termination signals to close connections and save the cache.
- **Thread Safety**: Uses mutex locks to ensure safe concurrent access to the cache.

## Installation

### Prerequisites

- Go 1.16 or higher

### Setup

1. Clone the repository:

   ```sh
   git clone https://github.com/Ndeta100/mcache.git
   go get go get github.com/Ndeta100/mcache@v1.0.1  # or v1.1.0
   cd mcache
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Build the server:

   ```sh
   go build -o cache-server
   ```

## Usage

### Starting the Server

Run the server executable:

```sh
./mcache
```

By default, the server listens on `localhost:6379`. You can modify the configuration (`config/config.yaml`) to change the listening address and port.

### Connecting to the Server

You can use a telnet client to connect to the server:

```sh
telnet localhost 6379
```

### Commands

- **SET**: Set a key-value pair in the cache.

  ```sh
  SET key value
  OK
  ```

- **GET**: Retrieve a value by key.

  ```sh
  GET key
  VALUE: value
  ```

- **DELETE**: Delete a key from the cache.

  ```sh
  DELETE key
  OK
  ```

- **Graceful Shutdown**: Press `CTRL+C` to stop the server. All active connections will be closed, and the cache state will be saved to disk.

## Configuration

The server configuration is defined in `config/config.yaml`. You can set options like:

- `host`: The server's hostname or IP.
- `port`: The server's port number.
- `capacity`: The maximum number of items allowed in the cache.
- `ttl`: The time-to-live for cache entries, in seconds.

## Graceful Shutdown

The server supports graceful shutdown, allowing it to:

- Stop accepting new connections.
- Close all active connections.
- Save the current cache state to disk.

This is triggered when the server receives a termination signal (`SIGINT` or `SIGTERM`).

## Development

### Running Tests

You can run unit tests for the cache and server:

```sh
go test ./...
```

### Code Structure

- **`server`**: Contains the TCP server code.
- **`store`**: Implements the key-value cache with persistence and TTL.
- **`handler`**: Handles client commands (`GET`, `SET`, `DELETE`).
- **`config`**: Contains configuration utilities.

## Improvements to Consider

- **Logging**: Improve logging for better observability of server events.
- **Authentication**: Add authentication to limit access to the cache.
- **REST API**: Implement a RESTful API for more accessible programmatic interaction.
- **Optimized Persistence**: Implement incremental snapshots or append-only logs for more efficient persistence.

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
