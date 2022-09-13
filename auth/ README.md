# Auth - authentication web service.

**Auth** - authentication service, developed for UNotes(notes system).

## API Methods

- `/api/oauth2/sign-in`
- `/api/oauth2/sign-up`
- `/api/oauth2/sign-out`
- `/api/oauth2/refresh`

## Run

1) You need to install [Docker](https://docs.docker.com/get-docker)
   and [Docker Compose](https://docs.docker.com/compose/install).
2) Open terminal.
3) `cd <folder with the project>/unotes/auth`
4) Setting up config [config.json](https://github.com/udholdenhed/unotes/blob/master/auth/configs/config.json).
5) `docker-compose up -d --build --remove-orphans`

## Development

### Prerequisites

- Recommended IDEs
    - [GoLand](https://www.jetbrains.com/go) (2022.2.2 and above).
    - [Visual Studio Code](https://code.visualstudio.com) (1.70 and above).

### Dependencies

1) `cd <folder with the project>/unotes/auth`
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

## Contributions

If you have **questions**, **ideas**, or you find a **bug**, you can create
an [issue,](https://github.com/udholdenhed/unotes/issues) and it will be reviewed. If you want to contribute to
the source code, fork this repository (`master`), realize your ideas and then create a new pull request. **Feel free!**
