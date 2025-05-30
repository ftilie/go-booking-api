# go-booking-api
Event Booking REST API powered by Go

This repository showcases a basic implementation of a REST API in Go, designed to manage event bookings efficiently. It includes features such as user authentication, event management, and booking operations, making it a great starting point for learning or building production-ready APIs.

## Features
- User authentication and authorization using secure mechanisms
- Full CRUD operations for managing events and bookings
- JSON-based API responses for seamless integration with front-end applications
- Middleware for logging, error handling, and request validation

## Prerequisites
- Go 1.18 or higher installed on your system

## Installation
1. Clone the repository:
    ```bash
    git clone https://github.com/ftilie/go-booking-api.git
    ```
2. Navigate to the project directory:
    ```bash
    cd go-booking-api
    ```
3. Install dependencies:
    ```bash
    go mod tidy
    ```
4. Run the application:
    ```bash
    go run main.go
    ```

## Testing
Run the test suite to ensure everything is working as expected:
```bash
go test ./...
```
