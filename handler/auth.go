package handler

import "net/http"

type AuthHandler struct {
}

func NewAuthHandler() AuthHandler {
	return AuthHandler{}
}

func (authHandler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login"))
}

func (authHandler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Register"))
}
