# AuthService - authentication web service.

**AuthService** - authentication service, developed for UNotes(notes system).

## API Methods

- `/api/auth/oauth2/sign-in`
- `/api/auth/oauth2/sign-up`
- `/api/auth/oauth2/sign-out`
- `/api/auth/oauth2/refresh`

## Run

### Docker

1) You need to install [Docker](https://docs.docker.com/get-docker).
2) In the root directory "auth-service" you need to create a new file named ".env".
3) In the ".env" file you need to add the following environment variables:
    - Auth Service:
        - AUTH_SERVICE_ACCESS_TOKEN_SECRET
        - AUTH_SERVICE_REFRESH_TOKEN_SECRET
    - PostgreSQL:
        - AUTH_SERVICE_POSTGRESQL_HOST
        - AUTH_SERVICE_POSTGRESQL_PORT
        - AUTH_SERVICE_POSTGRESQL_USERNAME
        - AUTH_SERVICE_POSTGRESQL_PASSWORD
        - AUTH_SERVICE_POSTGRESQL_DBNAME
        - AUTH_SERVICE_POSTGRESQL_SSLMODE
    - RedisDB:
        - AUTH_SERVICE_REDIS_ADDR
        - AUTH_SERVICE_REDIS_PASSWORD
        - AUTH_SERVICE_REDIS_DB
4) Now run the following commands:

- `docker build --tag auth-service .`
- `docker run --publish 8081:8081 --name auth-service --detach --restart always --env-file ./.env auth-service`

### Docker Compose

1) You need to install [Docker](https://docs.docker.com/get-docker).
2) And install [Docker Compose](https://docs.docker.com/compose/install).
3) In the root directory "auth-service" you need to create a new file named ".env".
4) In the ".env" file you need to add the following environment variables:
    - Auth Service:
        - AUTH_SERVICE_ACCESS_TOKEN_SECRET
        - AUTH_SERVICE_REFRESH_TOKEN_SECRET
    - PostgreSQL:
        - AUTH_SERVICE_POSTGRESQL_HOST
        - AUTH_SERVICE_POSTGRESQL_PORT
        - AUTH_SERVICE_POSTGRESQL_USERNAME
        - AUTH_SERVICE_POSTGRESQL_PASSWORD
        - AUTH_SERVICE_POSTGRESQL_DBNAME
        - AUTH_SERVICE_POSTGRESQL_SSLMODE
    - RedisDB:
        - AUTH_SERVICE_REDIS_ADDR
        - AUTH_SERVICE_REDIS_PASSWORD
        - AUTH_SERVICE_REDIS_DB
5) Now run the following command: `docker-compose up -d --build --remove-orphans`.

## Development

### Prerequisites

- Recommended IDEs
    - [GoLand](https://www.jetbrains.com/go) (2022.2.2 and above).
    - [Visual Studio Code](https://code.visualstudio.com) (1.70 and above).

### Dependencies

Run one of the following commands:

- `go mod download`
- `make go-dep`

### Environment

You must configure the environment variables:

- Auth Service:
    - AUTH_SERVICE_ACCESS_TOKEN_SECRET
    - AUTH_SERVICE_REFRESH_TOKEN_SECRET
- PostgreSQL:
    - AUTH_SERVICE_POSTGRESQL_HOST
    - AUTH_SERVICE_POSTGRESQL_PORT
    - AUTH_SERVICE_POSTGRESQL_USERNAME
    - AUTH_SERVICE_POSTGRESQL_PASSWORD
    - AUTH_SERVICE_POSTGRESQL_DBNAME
    - AUTH_SERVICE_POSTGRESQL_SSLMODE
- RedisDB:
    - AUTH_SERVICE_REDIS_ADDR
    - AUTH_SERVICE_REDIS_PASSWORD
    - AUTH_SERVICE_REDIS_DB

### Run

Run one of the following commands:

- `go run ./cmd/auth/`
- `go-run`

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
    - URL `http://localhost:8081/swagger/index.html#/`
    - If the swagger documentation has been changed, you can use the following command to generate a new one.
        - `swag init -g ./internal/handler/handler.go`

- If you want to change the config that is used (by default "production.env") you can use the environment variable
  named "ENVIRONMENT". To use "stage.env" use "ENVIRONMENT=STAGE", to use "development.env" use "
  ENVIRONMENT=DEVELOPMENT", for use in production, this variable may not be set.
