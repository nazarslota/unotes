.DEFAULT_GOAL := go-run

go-run: go-dep go-vet
	go run ./cmd/user/

go-build: go-dep go-vet
	go build -o ./build/ ./cmd/...

go-dep:
	go mod download

go-vet:
	go vet ./...

go-test:
	go test ./tests/...
