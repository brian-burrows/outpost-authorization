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

## Some helpful style guidelines

https://github.com/uber-go/guide/blob/master/style.md

## Authorization flows

### Bearer tokens

- Client sends a Token with each request.
- API verifies or rejects the token.

Standard approach in API design, because it's fast and stateless

#### Account creation Pseudocode

#### Login Pseudocode

- User sends a login request to `POST /login/`
  - Payload : `providerType: password, providerKey: 'username', credential: 'password'`
  - Server : Receives username, password
    - Hash password via bcrypt with secret key
    - Fetch `User` instance by `username`
    - Validate hashed password matches that in the database
    - Form accessToken JWT with `User` instance attributes, return to user.
    - Form refreshToken JWT with
  - Client : Stores JWT in sessions or cookies

#### Request Pseudocode

- User sends a request to any endpoint that requires authentication and authorization
  - Payload (in header): `Bearer: <JWT> token`
  - Server: Recieved request with header
    - Extracts bearer token
    - Decodes JWT, re-encrpyts the payload, check's that it matches JWT signature
  - If valid, processes the request

#### Refresh Pseudocode

- Client: Receives 401 Unauthorized.
- Client: Sends POST /refresh with the refreshToken in the body (or a secure cookie).
- Server: Decodes the Refresh JWT to get the userId and sessionId.
- Lookup: Finds the User in MongoDB and looks inside the sessions array.
- Match: Does the token sent by the client match one in the sessions array?
- Rotate: (Optional but Recommended) Delete the old session object and create a brand new one. This prevents "replay attacks."
- Respond: Return a new accessToken and the new refreshToken.

### OAuth2 + JWT

#### Login Pseudocode

- Server: Extracts sub (ID) and email.
  - Step A: User.findOne({ "providers.providerKey": sub }). If found, Log In.
  - Step B: If not found, User.findOne({ "email": email }).
    - If found: Link accounts. Push a new object into the providers array for Google.
    - If not found: Register. Create a new User document.
- Step C: Create a new entry in the sessions array and return the tokens.

#### Request Pseudocode

Same logic as a JWT flow, using internal token.

#### Refresh Pseudocode

Works exactly like the standard Refresh flow.

### Biometrics

#### Login Pseudocode

- The user's device generates a key pair.
  - The Private Key stays in the phone's Secure Enclave (FaceID/TouchID).
  - The Public Key is sent to your server.

#### Request Pseudocode

- Login: The server sends a "Challenge" (a random string).
- The phone signs that challenge using the Private Key.
- Your server verifies the signature using the public key

#### Refresh Pseudocode

## Example Schema

```json
{
  "_id": ObjectId("..."),
  "email": "user@example.com", // The unique "anchor" for the account
  "name": "Jane Doe",
  "sessions": [
    {
        "token": "abc...",
        "device": "iPhone 15",
        "expiresAt": ISODate("...")
    },
    {
        "token": "xyz...",
        "device": "Chrome / MacOS",
        "expiresAt": ISODate("...")
    }
  ]

  // We store an array of authentication methods
  "providers": [
    {
      "providerType": "password",
      "providerKey": "user@example.com",
      "credential": "$2b$12$Kls...", // Only used for password login
      "lastUsed": ISODate("2024-05-01..."),
      "isVerified": true,
    },
    {
      "providerType": "google",
      "providerKey": "1034567890123456789", // Google's 'sub' field (Subject ID)
      "credential": null,
      "lastUsed": ISODate("2024-05-10..."),
      "isVerified": true,
    }
  ],

  "role": "user",
  "createdAt": ISODate("...")
}
```
