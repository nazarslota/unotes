FROM golang:1.20-alpine AS builder

WORKDIR /go/src/github.com/nazarslota/unotes/note/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/ ./cmd/...

FROM alpine:latest

RUN apk --no-cache add ca-certificates

EXPOSE 8092
EXPOSE 8082

WORKDIR /root/

COPY --from=builder /go/src/github.com/nazarslota/unotes/note/configs/ ./configs/
COPY --from=builder /go/src/github.com/nazarslota/unotes/note/build/note ./

CMD ["./note"]
