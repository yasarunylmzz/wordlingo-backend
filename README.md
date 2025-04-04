# Project Name

## Description

This project is a Golang-based web application using the Echo framework. It connects to a PostgreSQL database, executes SQL schema migrations, and provides RESTful API endpoints.

## Features

- RESTful API built with Echo
- PostgreSQL database integration
- Automatic database creation and schema migration using a Makefile
- Configurable via `config.json`

## Prerequisites

Before running the project, ensure you have the following installed:

- [Go](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)
- [SQLC](https://sqlc.dev/)
- [Make](https://www.gnu.org/software/make/)

## Setup

### 1. Install Dependencies

Run the following command to install required Go modules:

```sh
go mod tidy
```

### 2. Database Initialization

Use the Makefile to create and configure the database.

- Clean and recreate the database:

  ```sh
  make all
  ```

  This executes:

  1. `make clean` - Drops the database if it exists.
  2. `make createdb` - Creates a new database.
  3. `make schema` - Applies the schema from `db/schema/schema.sql`.
  4. `make run_sql` - Runs all SQL files in `db/schema/`.

- If you only want to apply schema updates:
  ```sh
  make schema
  ```

## Running the Application

Start the Echo server using:

```sh
go run server.go
```

The server will start and listen on `localhost:1323`.

## Makefile Commands Explained

### `make all`

Runs `clean`, `createdb`, `schema`, and `run_sql` in sequence.

### `make clean`

Drops the existing database (if it exists).

### `make createdb`

Creates the database if it does not exist.

### `make schema`

Loads the schema from `db/schema/schema.sql`.

### `make run_sql`

Executes all SQL files in the `db/schema/` directory in order.

### `make help`

Displays a list of available `make` commands.

## API Endpoints

| Method | Endpoint             | Description                                        |
| ------ | -------------------- | -------------------------------------------------- |
| POST   | `/createuser`        | Create a new user                                  |
| GET    | `/loginuser`         | Login User endpoint                                |
| POST   | `/verificiationuser` | Verification Code with six-digit verification code |

## License

MIT License
