package main

import (
	"fmt"
	"go-19/database"
	"go-19/handler"
	"go-19/repository"
	"go-19/service"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	middlechi "github.com/go-chi/chi/v5/middleware"
)

func main() {
	var templates = template.Must(template.New("").ParseGlob("view/*.html"))
	db, err := database.InitDB()

	if err != nil {
		panic(err)
	}
	// assignmentRepo := repository.NewAssignmentRepository(db)
	// submissionRepo := repository.NewSubmissionRepo(db)
	// userRepo := repository.NewUserRepository(db)

	repo := repository.NewRepository(db)
	service := service.NewService(repo)

	// authService := service.NewAuthService(repo)
	// userService := service.NewUserService(repo)
	// submissionService := service.NewSubmissionService(repo)
	// assignmentService := service.NewAssignmentService(repo)

	assignmentHandler := handler.NewAssignmentHandler(service, templates)
	submissionHandler := handler.NewSubmissionHandler(submissionService, userService, assignmentService, templates)
	authHandler := handler.NewAuthHandler(authService, templates)

	r := chi.NewRouter()
	r.Use(middlechi.Logger)

	r.Post("/login", authHandler.HandleLogin)

	r.Get("/student/home", assignmentHandler.ListAssignments)
	r.Get("/student/submit", assignmentHandler.ShowSubmitForm)
	r.Post("/student/submit", assignmentHandler.SubmitAssignment)

	r.Get("/lecturer/home", submissionHandler.Home)
	r.Get("/lecturer/grade-form", submissionHandler.ShowGradeForm)
	r.Post("/lecturer/grade", submissionHandler.GradeSubmission)
	r.Get("/", authHandler.ShowLoginForm)

	fs := http.FileServer(http.Dir("view"))
	r.Handle("/view/*", http.StripPrefix("/view/", fs))

	fmt.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("error server")
	}
}
