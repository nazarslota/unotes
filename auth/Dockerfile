FROM golang:1.19.1-alpine AS builder

WORKDIR /go/src/github.com/udholdenhed/unotes/auth-service/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/ ./cmd/...

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

EXPOSE 8081, 8091

COPY --from=builder /go/src/github.com/udholdenhed/unotes/auth-service/configs/ ./configs/
COPY --from=builder /go/src/github.com/udholdenhed/unotes/auth-service/build/auth ./

CMD ["./auth"]
