# URL Shortening Service (Backend)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This project is a backend service for shortening URLs. It provides a simple RESTful API to create short URLs and return the original URLs. It is based on the [URL Shortening Service](https://roadmap.sh/projects/url-shortening-service) project from [Roadmap.sh](https://roadmap.sh/), but with some additional features, authentications, and customization.

## Features ✨

- User authentication with JSON Web Token (JWT).
- Shorten long URLs.
- Return the original URLs using the short code.
- Edit the original URLs.
- Delete the short URLs.
- Track the number of clicks on each short URL.
- Display statistics of the short URL: last accessed time, total clicks last 7 days, 30 days, and 90 days (per day).

## Technologies Used ⚙️

- **Programming Language**: Golang.
- **DBMS**: PostgreSQL.

## Installation 🛠️

1. Clone the repository.

```sh
git clone https://github.com/azbagas/url-shortening-backend.git
```

2. Navigate to the project directory.

```sh
cd url-shortening-backend
```

3. Download dependencies.

```sh
go mod tidy
```

## Usage 📖

1. Copy the `.env.example` file into `.env` and fill in the values with your own configuration.

```sh
cp .env.example .env
```

2. Create a PostgreSQL database with the name you have specified in the `.env` file.

3. Run the database migrations.

> Note: The migration tool used is Golang Migrate. If you haven't installed it yet, check their documentation [here](https://github.com/golang-migrate/migrate), or follow the optional command below:

```sh
# Optional: Install Golang Migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Migrate
migrate -database "postgres://postgres:password@localhost:5432/<your_database_name>?sslmode=disable" -path db/migrations up
```

4. Start the server.

> Note: Since this app is using [Google Wire](https://github.com/google/wire) for dependency injection, we need to run two files at once: `wire_gen.go` and `main.go`

```sh
go run wire_gen.go main.go
```

## API Endpoints 📃

List of API endpoints can be accessed in `api_doc.yml` file.

## Running Tests 🧪

To run tests, run the following command:

```bash
  go test -v ./test
```

> ⚠️ **Warning**: The test will truncate all the database tables. So, don't run it in a production environment.
