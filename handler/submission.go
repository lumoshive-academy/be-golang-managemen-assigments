// handler/lecturer_handler.go
package handler

import (
	"go-19/service"
	"html/template"
	"net/http"
	"strconv"
)

type SubmissionHandler struct {
	Submissionservice service.SubmissionService
	Userservice       service.UserService
	AssignmentService service.AssignmentService
	template          *template.Template
}

func NewSubmissionHandler(submissionService service.SubmissionService, userService service.UserService, assignmentService service.AssignmentService, tmpl *template.Template) *SubmissionHandler {
	return &SubmissionHandler{
		Submissionservice: submissionService,
		Userservice:       userService,
		AssignmentService: assignmentService,
		template:          tmpl,
	}
}

func (h *SubmissionHandler) Home(w http.ResponseWriter, r *http.Request) {
	submissions, err := h.Submissionservice.GetAllSubmissions()
	if err != nil {
		http.Error(w, "Gagal mengambil data submission", http.StatusInternalServerError)
		return
	}

	h.template.ExecuteTemplate(w, "lecturer_home", submissions)
}

func (h *SubmissionHandler) ShowGradeForm(w http.ResponseWriter, r *http.Request) {
	studentIDStr := r.URL.Query().Get("student_id")
	assignmentIDStr := r.URL.Query().Get("assignment_id")

	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		http.Error(w, "Invalid student_id", http.StatusBadRequest)
		return
	}

	assignmentID, err := strconv.Atoi(assignmentIDStr)
	if err != nil {
		http.Error(w, "Invalid assignment_id", http.StatusBadRequest)
		return
	}

	// Ambil data untuk ditampilkan di form
	student, err := h.Userservice.GetUserByID(studentID)
	if err != nil {
		http.Error(w, "Student not found", http.StatusInternalServerError)
		return
	}

	assignment, err := h.AssignmentService.GetAssignmentByID(assignmentID)
	if err != nil {
		http.Error(w, "Assignment not found", http.StatusInternalServerError)
		return
	}

	data := struct {
		StudentID       int
		AssignmentID    int
		StudentName     string
		AssignmentTitle string
	}{
		StudentID:       student.ID,
		AssignmentID:    assignment.ID,
		StudentName:     student.Name,
		AssignmentTitle: assignment.Title,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.template.ExecuteTemplate(w, "grade_form", data)
}

func (h *SubmissionHandler) GradeSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Gagal parsing form", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(r.FormValue("student_id"))
	if err != nil {
		http.Error(w, "Invalid student_id", http.StatusBadRequest)
		return
	}

	assignmentID, err := strconv.Atoi(r.FormValue("assignment_id"))
	if err != nil {
		http.Error(w, "Invalid assignment_id", http.StatusBadRequest)
		return
	}

	gradeStr := r.FormValue("grade")
	grade, err := strconv.ParseFloat(gradeStr, 64)
	if err != nil {
		http.Error(w, "Invalid grade", http.StatusBadRequest)
		return
	}

	err = h.Submissionservice.GradeSubmission(studentID, assignmentID, grade)
	if err != nil {
		http.Error(w, "Gagal memberi nilai: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/lecturer/home", http.StatusSeeOther)
}
