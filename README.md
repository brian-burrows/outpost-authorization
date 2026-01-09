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

## Testing, Test Coverage, and Complexity Scores

### Viewing test coverage

```bash
go test -coverprofile=cp.out ./...
go tool cover -html=cp.out
```

### Viewing code complexity

This script calculates the ABC score (assignments, branches, conditionals) of a file.

```bath
go build -o complexity ./cmd/complexity
./complexity ./<filepath>.go
```

For a single function, the scores roughly indicate:

- **Simple code: (ideal)** 0 - 5
- **Fairly Simple: (OK)** 6 - 10
- **Moderatily Complex** 11 - 15
- **High complexity** 16 - 20
- **Very high complexity** > 20

Moderately complex code should be reviewed.
High complexity code is a candidate for a refactor.
Very high complexity code should be decomposed.

Note: This looks at function-level complexity.
It cannot distinguish between

- One giant function (high ABC score)
- 20 trivially tiny functions (all low ABC score)
- 5 well-designed, cohesive fnuctions (all moderately low ABC scores)

So, it can be useful to aggregate these.

- Total ABC per package
- Avg ABC per function
- Max ABC
- Function count
