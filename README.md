# Go Task Manager

A simple task management application built with Go, featuring a REST API server and a CLI client. This project demonstrates core Go concepts including structs, interfaces, error handling, concurrency, HTTP programming, and file I/O.

## Project Structure

```
taskmanager/
├── cmd/
│   ├── client/        # CLI client application
│   │   └── main.go
│   └── server/        # HTTP server application
│       └── main.go
├── internal/
│   ├── api/           # API handlers
│   │   └── handler.go
│   ├── models/        # Data models
│   │   └── task.go
│   └── storage/       # Data persistence
│       └── file_storage.go
├── go.mod             # Go module file
├── go.sum             # Go module checksum
└── README.md          # This file
```

## Features

- **Task Management**: Create, read, update, and delete tasks
- **REST API**: HTTP endpoints for task operations
- **CLI Client**: Command-line interface to interact with the API
- **Persistent Storage**: File-based JSON storage
- **Concurrency-Safe**: Thread-safe operations with mutex locks

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/NanobyteRuata/taskmanager.git
   cd taskmanager
   ```

2. Install dependencies:
   ```
   go mod download
   ```

## Usage

### Starting the Server

```bash
go run cmd/server/main.go
```

The server will start on port 8080 by default. You can change the port by setting the `PORT` environment variable.

### Using the CLI Client

The CLI client provides several commands to interact with the task manager:

#### List all tasks
```bash
go run cmd/client/main.go list
```

#### Add a new task
```bash
go run cmd/client/main.go add "Task description" [due-date]
```
Example:
```bash
go run cmd/client/main.go add "Learn Go" 2023-12-31
```

#### Mark a task as completed
```bash
go run cmd/client/main.go complete [task-id]
```

#### Delete a task
```bash
go run cmd/client/main.go delete [task-id]
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /tasks | List all tasks |
| POST | /tasks | Create a new task |
| GET | /tasks/{id} | Get a specific task |
| PATCH | /tasks/{id}/complete | Mark a task as completed |
| DELETE | /tasks/{id} | Delete a task |

## Data Model

A task has the following structure:

```go
type Task struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
    CompletedAt time.Time `json:"completed_at,omitempty"`
    DueDate     time.Time `json:"due_date,omitempty"`
}
```

## Go Concepts Demonstrated

- **Structs and Methods**: Task struct with associated methods
- **Interfaces**: TaskStore interface for storage abstraction
- **Error Handling**: Custom error types and proper error propagation
- **Concurrency**: Goroutines, channels, and mutex for thread safety
- **HTTP Programming**: RESTful API with proper status codes and JSON responses
- **JSON Handling**: Marshaling and unmarshaling of JSON data
- **File I/O**: Reading from and writing to files
- **Command-Line Arguments**: Parsing and handling CLI arguments
- **Graceful Shutdown**: Proper server shutdown with signal handling

## Future Enhancements

- Add authentication and user management
- Implement database storage (e.g., SQLite, PostgreSQL)
- Add task categories and priorities
- Create a web UI
- Add unit and integration tests

## License

[MIT License](LICENSE)
