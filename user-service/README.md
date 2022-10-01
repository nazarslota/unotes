# UserService - service for saving user data.

**UserService** - service for saving user data, developed for UNotes(notes system).

## API Methods

- `/api/user/create`
- `/api/user/find-id`
- `/api/user/find-username`

## Run

1) You need to install [Docker](https://docs.docker.com/get-docker)
   and [Docker Compose](https://docs.docker.com/compose/install).
2) Open terminal.
3) `cd /unotes/auth-service/`
4) Create an `.env` file and set the necessary environment variables.
5) `docker-compose up -d --build --remove-orphans`

## Development

### Prerequisites

- Recommended IDEs
    - [GoLand](https://www.jetbrains.com/go) (2022.2.2 and above).
    - [Visual Studio Code](https://code.visualstudio.com) (1.70 and above).

### Dependencies

1) `cd <folder with the project>/unotes/auth-service`
2) `go mod download` or `make go-dep`

### Other

- For easy management of database migrations, it is recommended to use
  the [migrate](https://github.com/golang-migrate/migrate) tool.
    - Up
        - `migrate -path ./schema -database 'postgres://<username>:<password>@<host>:<port>/<database>' up`
        - `migrate -path ./schema -database 'postgres://<username>:<password>@<host>:<port>/<database>?sslmode=disable' up`
    - Down
        - `migrate -path ./schema -database 'postgres://<username>:<password>@<host>:<port>/<database>' down`
        - `migrate -path ./schema -database 'postgres://<username>:<password>@<host>:<port>/<database>?sslmode=disable' down`

- The project also uses swagger documentation.
    - URL `http://localhost:8082/swagger/index.html#/`
    - If the swagger documentation has been changed, you can use the following command to generate a new one.
        - `swag init -g ./internal/handler/handler.go`
