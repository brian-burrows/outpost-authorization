FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN echo 'package main\nimport ("fmt"; "time")\nfunc main() { fmt.Println("Dummy Service Running"); for { time.Sleep(1 * time.Hour) } }' > dummy.go

RUN go build -o /dummy_service dummy.go

FROM alpine:latest

COPY --from=builder /dummy_service /dummy_service

ENTRYPOINT ["/dummy_service"]