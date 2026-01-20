FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN <<EOF cat > dummy.go
package main
import (
    "fmt"
    "time"
)
func main() {
    fmt.Println("Dummy Service Running")
    for {
        time.Sleep(1 * time.Hour)
    }
}
EOF

RUN go build -o /dummy_service dummy.go

FROM alpine:latest
COPY --from=builder /dummy_service /dummy_service
ENTRYPOINT ["/dummy_service"]