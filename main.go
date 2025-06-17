package main

import (
	"fmt"
	"go-19/handler"
	"go-19/middleware"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	middlechi "github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middlechi.Logger)
	r.Use(middleware.Auth)
	authHandler := handler.NewAuthHandler()

	r.Get("/login-view", authHandler.Login)
	r.Post("/Register", authHandler.Register)
	r.Get("/users/{userID}/posts/{postID}", func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		postID := chi.URLParam(r, "postID")
		w.Write([]byte("User ID: " + userID + ", Post ID: " + postID))
	})

	fmt.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("error server")
	}

}
