# Endpoints design

## Application Programming Interface (API)

- `CreateUser(ProviderType, ProviderKey, Credential)`

### Domain Model

#### User

```go
type User struct {
    ID          string
    Email       string
    Roles       []Role
    Permissions []Permission
    CreatedAt   time.Time
}
```
