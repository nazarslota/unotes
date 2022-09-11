.DEFAULT_GOAL := go-run

go-run: go-dep go-vet go-test
	go run ./cmd/auth/

go-build: go-dep go-vet go-test
	go build -o ./build/ ./cmd/...

go-dep:
	go mod download

go-vet:
	go vet ./...

go-test:
	go test ./tests/...
