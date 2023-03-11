# Auth

This is an authentication and authorization microservice written in [Go](https://github.com/golang/go) and built using
REST and [gRPC](https://grpc.io). It uses PostgreSQL and Redis for storage and implements the OAuth2.0
authentication method for secure authentication.

## Run

### Docker

1. You must have [Docker](https://docs.docker.com/engine/install) installed.
2. Now in the main folder of the project you need to create an `.env` file in which you need to configure to add the
   following environment variables.
   ```
   AUTH_ACCESS_TOKEN_SECRET=
   AUTH_REFRESH_TOKEN_SECRET=
   
   AUTH_POSTGRESQL_HOST=
   AUTH_POSTGRESQL_PORT=
   AUTH_POSTGRESQL_USERNAME=
   AUTH_POSTGRESQL_PASSWORD=
   AUTH_POSTGRESQL_DBNAME=
   AUTH_POSTGRESQL_SSLMODE=
   
   AUTH_REDIS_ADDR=
   AUTH_REDIS_PASSWORD=
   AUTH_REDIS_DB=
   ```
3. Now run the following commands to build and run the Docker container.
   ```
   docker build --tag auth .
   docker run --publish 8082:8082 --publish 8092:8092 --name auth --detach --restart always --env-file ./.env auth
   ```

### Docker Compose

1. You must have [Docker](https://docs.docker.com/engine/install)
   and [Docker Compose](https://docs.docker.com/compose/install) installed.
2. Now in the main folder of the project you need to create an `.env` file in which you need to configure to add the
   following environment variables.
   ```
   AUTH_ACCESS_TOKEN_SECRET=
   AUTH_REFRESH_TOKEN_SECRET=
   
   AUTH_POSTGRESQL_HOST=
   AUTH_POSTGRESQL_PORT=
   AUTH_POSTGRESQL_USERNAME=
   AUTH_POSTGRESQL_PASSWORD=
   AUTH_POSTGRESQL_DBNAME=
   AUTH_POSTGRESQL_SSLMODE=
   
   AUTH_REDIS_ADDR=
   AUTH_REDIS_PASSWORD=
   AUTH_REDIS_DB=
   ```
3. Now you need to run the following command.
   ```
   docker-compose up --detach --build --remove-orphans
   ```

## Development

### Prerequisites

* [GoLand](https://www.jetbrains.com/go) or [Visual Studio Code](https://code.visualstudio.com)
* [Go](https://go.dev/dl)
* [Protocol Buffer Compiler](https://grpc.io/docs/protoc-installation)

### Dependencies

```
make dependencies
```

### Other

#### [Protocol Buffer](https://grpc.io/docs/languages/go/quickstart)

- Code generation.
    ```
    make protobuf
    ```

- Plugins install.
    ```
    make protoplugins
    ```

- Plugins.
    - [protoc-gen-validate](https://github.com/bufbuild/protoc-gen-validate)

#### [Swagger](https://swagger.io)

- URL: http://localhost:8081/api/swagger/index.html

- Code generation.
    ```
    make swagger
    ```

- Tools install.
    ```
    make tools
    ```

- Tools.
    - [swag](https://github.com/swaggo/swag)

#### Migrations

- Using [migrate](https://github.com/golang-migrate/migrate) tool.
    - Up.
        ````
        migrate -path ./schema -database 'postgres://<username>:<password>@<host>:<port>/<database>' up
        migrate -path ./schema -database 'postgres://<username>:<password>@<host>:<port>/<database>?sslmode=disable' up
        ````

    - Down.
        ````
        migrate -path ./schema -database 'postgres://<username>:<password>@<host>:<port>/<database>' down
        migrate -path ./schema -database 'postgres://<username>:<password>@<host>:<port>/<database>?sslmode=disable' down`
        ````

- Tools install.
    ```
    make tools
    ```

- Tools.
    - [migrate](https://github.com/golang-migrate/migrate)

#### Environment variables

```
AUTH_ACCESS_TOKEN_SECRET=
AUTH_REFRESH_TOKEN_SECRET=

AUTH_POSTGRESQL_HOST=
AUTH_POSTGRESQL_PORT=
AUTH_POSTGRESQL_USERNAME=
AUTH_POSTGRESQL_PASSWORD=
AUTH_POSTGRESQL_DBNAME=
AUTH_POSTGRESQL_SSLMODE=
   
AUTH_REDIS_ADDR=
AUTH_REDIS_PASSWORD=
AUTH_REDIS_DB=
```
