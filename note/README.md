# Note

This microservice is written in [Go](https://github.com/golang/go) and implements a [gRPC](https://grpc.io) service,
utilizing [gRPC Gateway](https://github.com/grpc-ecosystem/grpc-gateway) for HTTP/REST endpoint compatibility
and [Swagger](https://swagger.io/) for API documentation. It also integrates with [MongoDB](https://www.mongodb.com)
for data storage and [Docker](https://www.docker.com) for easy deployment.

## Run

### Docker

1. You must have [Docker](https://docs.docker.com/engine/install) installed.
2. Now in the main folder of the project you need to create an `.env` file in which you need to configure to add the
   following environment variables.
   ```
   NOTE_MONGODB_HOST=
   NOTE_MONGODB_PORT=
   NOTE_MONGODB_USERNAME=
   NOTE_MONGODB_PASSWORD=
   NOTE_MONGODB_DATABASE=
   ```
3. Now run the following commands to build and run the Docker container.
   ```
   docker build --tag note .
   docker run --publish 8082:8082 --publish 8092:8092 --name note --detach --restart always --env-file ./.env note
   ```

### Docker Compose

1. You must have [Docker](https://docs.docker.com/engine/install)
   and [Docker Compose](https://docs.docker.com/compose/install) installed.
2. Now in the main folder of the project you need to create an `.env` file in which you need to configure to add the
   following environment variables.
   ```
   NOTE_MONGODB_HOST=
   NOTE_MONGODB_PORT=
   NOTE_MONGODB_USERNAME=
   NOTE_MONGODB_PASSWORD=
   NOTE_MONGODB_DATABASE=
   ```
3. Now you need to run the following command.
   ```
   docker-compose up --detach --build --remove-orphans
   ```
