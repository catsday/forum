package handlers

import (
	"database/sql"
	"forum/internal/models"
	"html/template"
	"log"
	"net/http"
)

type ErrorData struct {
	Code        int
	Status      string
	Description string
}

func RenderError(w http.ResponseWriter, code int, description string) {
	w.WriteHeader(code)
	data := ErrorData{
		Code:        code,
		Status:      http.StatusText(code),
		Description: description,
	}

	files := []string{
		"./ui/templates/error.html",
		"./ui/templates/header.html",
		"./ui/templates/footer.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("Error rendering error page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Printf("Error executing error template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func AuthorizeAndHandle(db *sql.DB, handlerFunc func(w http.ResponseWriter, r *http.Request, userID int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userModel := &models.UserModel{DB: db}
		userID, err := userModel.GetSessionUserIDFromRequest(r)
		if err != nil {
			RenderError(w, http.StatusUnauthorized, "Access denied. Please log in to view this page.")
			return
		}

		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			RenderError(w, http.StatusMethodNotAllowed, "This resource does not support the HTTP method used.")
			return
		}

		handlerFunc(w, r, userID)
	}
}
