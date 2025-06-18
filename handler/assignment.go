package handler

import (
	"go-19/model"
	"go-19/service"
	"html/template"
	"net/http"
	"strconv"
)

type AssignmentHandler struct {
	Service  service.Service
	Template *template.Template
}

func NewAssignmentHandler(server service.Service, template *template.Template) AssignmentHandler {
	return AssignmentHandler{
		Service:  server,
		Template: template,
	}
}

func (assignmentHandler *AssignmentHandler) ListAssignments(w http.ResponseWriter, r *http.Request) {
	// Ambil user_id dari cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Konversi user_id ke int
	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Ambil data assignments dari service
	assignments, err := assignmentHandler.Service.AssignmentService.GetAllAssignments()
	if err != nil {
		http.Error(w, "Failed to fetch assignments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ambil data user (student) untuk menampilkan namanya
	student, err := assignmentHandler.Service.UserService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch student data", http.StatusInternalServerError)
		return
	}

	// Kirim ke template
	data := struct {
		StudentName string
		Assignments []model.Assignment
	}{
		StudentName: student.Name,
		Assignments: assignments,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := assignmentHandler.Template.ExecuteTemplate(w, "assignment_list", data); err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (assignmentHandler *AssignmentHandler) SubmitAssignment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form for file upload (max 10MB)
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Gagal membaca form: "+err.Error(), http.StatusBadRequest)
			return
		}

		assignmentID, err := strconv.Atoi(r.FormValue("assignment_id"))
		if err != nil {
			http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
			return
		}

		studentID, err := strconv.Atoi(r.FormValue("student_id"))
		if err != nil {
			http.Error(w, "Invalid student ID", http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "File tidak valid: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		status, err := assignmentHandler.Service.AssignmentService.SubmitAssignment(studentID, assignmentID, file, fileHeader)
		if err != nil {
			http.Error(w, "Gagal submit: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Berhasil submit dengan status: " + status))
	}
}

func (h *AssignmentHandler) ShowSubmitForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/student/home", http.StatusSeeOther)
		return
	}

	assignmentIDStr := r.URL.Query().Get("assignment_id")
	assignmentID, err := strconv.Atoi(assignmentIDStr)
	if err != nil {
		http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
		return
	}

	assignment, err := h.Service.AssignmentService.GetAssignmentByID(assignmentID)
	if err != nil {
		http.Error(w, "Assignment not found", http.StatusNotFound)
		return
	}

	// Ambil user_id dari cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.Service.UserService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	data := struct {
		Assignment  model.Assignment
		StudentID   int
		StudentName string
	}{
		Assignment:  *assignment,
		StudentID:   user.ID,
		StudentName: user.Name,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.Template.ExecuteTemplate(w, "submit_form", data)
}
