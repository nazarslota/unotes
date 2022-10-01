FROM golang:1.19.1-alpine AS builder

WORKDIR /go/src/github.com/udholdenhed/unotes/user-service/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/ ./cmd/...

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/src/github.com/udholdenhed/unotes/user-service/configs/ ./configs/
COPY --from=builder /go/src/github.com/udholdenhed/unotes/user-service/build/user ./

CMD ["./user"]
