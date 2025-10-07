# LSCS Core API

## Project Overview

This project is the official Core API microservice for the La Salle Computer Society (LSCS). It is a Go-based API that provides authentication and member data services. It uses `chi` for routing, `sqlc` for database queries, a MySQL database, and JWT for authentication.

## Building and Running

### Prerequisites

- Go
- MySQL

### Instructions

1.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

2.  **Set up Environment Variables:**
    Create a `.env` file in the root directory, using `.env.example` as a template. You will need to add a `JWT_SECRET` to this file.

3.  **Run the Database:**
    Make sure you have a running MySQL instance and the database schema is created using `schema.sql`.

4.  **Build and Run the Application:**
    The following commands can be used to build and run the application:

    *   **Build:**
        ```bash
        make build
        ```
        or
        ```bash
        go build -o bin/lscs-core-api ./cmd/api/main.go
        ```

    *   **Run:**
        ```bash
        make run
        ```
        or
        ```bash
        ./bin/lscs-core-api
        ```

    *   **Test:**
        ```bash
        make test
        ```
        or
        ```bash
        go test -v ./...
        ```

## Development Conventions

*   **API:** The API endpoints are defined in the `internal/server/routes.go` file.
*   **Authentication:** The API uses JWT for authentication. The `/auth/google/callback` endpoint is used to obtain a JWT.
*   **Database:** The database schema is defined in `schema.sql`. `sqlc` is used to generate type-safe Go code from SQL queries in `query.sql`. The configuration for `sqlc` is in `sqlc.yaml`.
*   **Routing:** The server uses `chi` for routing. Routes are registered in `internal/server/routes.go`.
*   **Main Entrypoint:** The main entrypoint of the application is `cmd/api/main.go`.
*   **Server:** The server is initialized in `internal/server/server.go`.