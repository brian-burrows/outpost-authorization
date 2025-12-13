# Outpost Authorization

This is intended to be an authorization microservice.

## Resources for Golang

- **Introduction to Go:** `https://go.dev/tour/welcome/1`
- **How to write Go Code:** `https://go.dev/doc/code`
- **Idiomatic Code:** ` `
- **Standard Library:** `https://pkg.go.dev/std`
- **In depth documentation:** `https://go.dev/doc/#articles`
- **Cheat sheet:** `https://github.com/a8m/golang-cheat-sheet`

## Module and package structure

## Getting started

Install dependencies

```bash
go mod tidy
```

You can run the main executable using `go run ./cmd/auth`.

You can build the standalone binary file to deploy (e.g., to a Docker container or server) by using the `go build` command:

```bash
# 1. Compile the program
go build ./cmd/auth
# By default, this creates an executable file in your current directory
# named after the directory it compiled: 'auth'.
# 2. Run the executable
# Note: You now run the binary directly, not using the 'go' tool
./auth
```

Or, you can output the file at a specific location with a specified name:

```bash
# Creates the executable file at ./bin/outpost-auth
go build -o ./bin/outpost-auth ./cmd/auth
# Run the compiled binary
./bin/outpost-auth
```
