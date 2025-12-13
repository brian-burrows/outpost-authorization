package auth

import "net/http"

var HandlerFileName string = "handler.go"

type UserHandlers struct {
	read_db  UserRepository
	write_db UserRepository
}

func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {}
func (h *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request)    {}

type TokenHandlers struct {
	tokenService TokenService
}

func (t *TokenHandlers) IssueToken(w http.ResponseWriter, r *http.Request)        {}
func (t *TokenHandlers) ValidateToken(w http.ResponseWriter, r *http.Request)     {}
func (t *TokenHandlers) IssueRefreshToken(w http.ResponseWriter, r *http.Request) {}
