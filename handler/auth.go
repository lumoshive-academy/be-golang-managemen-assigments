package handler

import (
	"go-19/service"
	"html/template"
	"net/http"
	"strconv"
)

type AuthHandler struct {
	authService service.AuthService
	template    *template.Template
}

func NewAuthHandler(authService service.AuthService, tmpl *template.Template) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		template:    tmpl,
	}
}

func (h *AuthHandler) ShowLoginForm(w http.ResponseWriter, r *http.Request) {
	h.template.ExecuteTemplate(w, "login", nil)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.authService.Login(email, password)
	if err != nil {
		http.Error(w, "Email atau password salah", http.StatusUnauthorized)
		return
	}

	// Simpan user_id dan role ke cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: strconv.Itoa(user.ID),
		Path:  "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "user_role",
		Value: user.Role,
		Path:  "/",
	})

	// Redirect berdasarkan role
	if user.Role == "lecturer" {
		http.Redirect(w, r, "/lecturer/home", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/student/home", http.StatusSeeOther)
	}
}
