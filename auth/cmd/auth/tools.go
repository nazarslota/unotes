//go:build tools
// +build tools

package main

import (
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/swaggo/swag/cmd/swag"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
