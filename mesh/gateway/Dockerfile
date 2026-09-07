FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git protobuf-dev

COPY go.mod go.sum ./
RUN go mod download

COPY proto/ /app/proto/
COPY main.go /app/

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN protoc --go_out=. --go-grpc_out=. proto/core.proto

RUN go build -o gateway main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gateway .

EXPOSE 8080

CMD ["./gateway"]
