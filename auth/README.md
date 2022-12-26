# AuthService - authentication web service.

**Auth** - authentication service, developed for UNotes(notes system).

## Run

### Docker

1) You need to install [Docker](https://docs.docker.com/get-docker).

2) In the root directory "auth" you need to create a new file named ".env".

3) In the ".env" file you need to add the following environment variables:

   Auth

    ````
    AUTH_SERVICE_ACCESS_TOKEN_SECRET
    AUTH_SERVICE_REFRESH_TOKEN_SECRET
    ````
   PostgreSQL

   ````
    AUTH_SERVICE_POSTGRESQL_HOST
    AUTH_SERVICE_POSTGRESQL_PORT
    AUTH_SERVICE_POSTGRESQL_USERNAME
    AUTH_SERVICE_POSTGRESQL_PASSWORD
    AUTH_SERVICE_POSTGRESQL_DBNAME
    AUTH_SERVICE_POSTGRESQL_SSLMODE
    ````
   RedisDB

    ````
    AUTH_SERVICE_REDIS_ADDR
    AUTH_SERVICE_REDIS_PASSWORD
    AUTH_SERVICE_REDIS_DB
    ````

4) Now run the following commands.
    ```
    docker build --tag auth .
    docker run --publish 8081:8081 --name auth --detach --restart always --env-file ./.env auth
    ```

### Docker Compose

1) You need to install [Docker](https://docs.docker.com/get-docker)
   and [Docker Compose](https://docs.docker.com/compose/install).

2) In the root directory "auth" you need to create a new file named ".env".

3) In the ".env" file you need to add the following environment variables.

   Auth

    ````
    AUTH_SERVICE_ACCESS_TOKEN_SECRET
    AUTH_SERVICE_REFRESH_TOKEN_SECRET
    ````
   PostgreSQL

   ````
    AUTH_SERVICE_POSTGRESQL_HOST
    AUTH_SERVICE_POSTGRESQL_PORT
    AUTH_SERVICE_POSTGRESQL_USERNAME
    AUTH_SERVICE_POSTGRESQL_PASSWORD
    AUTH_SERVICE_POSTGRESQL_DBNAME
    AUTH_SERVICE_POSTGRESQL_SSLMODE
    ````
   RedisDB

    ````
    AUTH_SERVICE_REDIS_ADDR
    AUTH_SERVICE_REDIS_PASSWORD
    AUTH_SERVICE_REDIS_DB
    ````

4) Now run the following command.
   ````
   docker-compose up -d --build --remove-orphans
   ````

## Development

### Prerequisites

- [GoLang](https://go.dev/dl) (1.19.4 recommended).
- [Protocol Buffer Compiler](https://grpc.io/docs/protoc-installation) (v3.15.8 recommended).

- Tools

    - [Swag](https://github.com/swaggo/swag) (1.8.9 recommended).
    - [Migrate](https://github.com/golang-migrate/migrate) (4.15.2 recommended).

- Recommended IDEs

    - [GoLand](https://www.jetbrains.com/go) (2022.2.2 and above).
    - [Visual Studio Code](https://code.visualstudio.com) (1.70 and above).

### Dependencies

Run the following command

````
make dependencies
````

### Environment

Auth

````
AUTH_SERVICE_ACCESS_TOKEN_SECRET
AUTH_SERVICE_REFRESH_TOKEN_SECRET
````

PostgreSQL

````
AUTH_SERVICE_POSTGRESQL_HOST
AUTH_SERVICE_POSTGRESQL_PORT
AUTH_SERVICE_POSTGRESQL_USERNAME
AUTH_SERVICE_POSTGRESQL_PASSWORD
AUTH_SERVICE_POSTGRESQL_DBNAME
AUTH_SERVICE_POSTGRESQL_SSLMODE
````

RedisDB

````
AUTH_SERVICE_REDIS_ADDR
AUTH_SERVICE_REDIS_PASSWORD
AUTH_SERVICE_REDIS_DB
````

### Run

Run the following commands.

````
make build
./bin/auth
````

### Other

#### For easy management of database migrations, it is recommended to use the [Migrate](https://github.com/golang-migrate/migrate) tool.

- Up
    - `migrate -path ./schema -database 'postgres://<username>:<password>@<host>:<port>/<database>' up`
    - `migrate -path ./schema -database 'postgres://<username>:<password>@<host>:<port>/<database>?sslmode=disable' up`
- Down
    - `migrate -path ./schema -database 'postgres://<username>:<password>@<host>:<port>/<database>' down`
    - `migrate -path ./schema -database 'postgres://<username>:<password>@<host>:<port>/<database>?sslmode=disable' down`

#### The project also uses swagger documentation. [Swag](https://github.com/swaggo/swag).

- URL [http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html)
- If the swagger documentation has been changed, you can use the following command to generate a new
  one - `make swagger`

#### In order to generate proto files, just run the following command. [Protocol Buffer Compiler](https://grpc.io/docs/protoc-installation) required.

- `.proto` and `.pb.go` files are stored in the `/api/proto` folder.
- If the proto files have been changed, you can use the following command - `make proto`. It will generate files in
  `/api/proto folder`.

#### Environment

If you want to change the config that is used (by default `production.env`) you can use the environment variable
named `ENVIRONMENT`. To use "stage.env" use `ENVIRONMENT=STAGE`, to use `development.env` use `ENVIRONMENT=DEVELOPMENT`,
for use in production, this variable may not be set.
