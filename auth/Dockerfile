FROM golang:1.19.4-alpine AS builder

WORKDIR /go/src/github.com/nazarslota/unotes/auth/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/ ./cmd/...

FROM alpine:latest

RUN apk --no-cache add ca-certificates

EXPOSE 8081
EXPOSE 8091

WORKDIR /root/

COPY --from=builder /go/src/github.com/nazarslota/unotes/auth/configs/ ./configs/
COPY --from=builder /go/src/github.com/nazarslota/unotes/auth/build/auth ./

CMD ["./auth"]
