package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type ErrorData struct {
	Code    int
	Message string
}

func RenderError(w http.ResponseWriter, code int, message string) {
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		log.Printf("Ошибка парсинга шаблона: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)

	data := ErrorData{
		Code:    code,
		Message: message,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Ошибка рендеринга шаблона: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ErrorBadRequest(w http.ResponseWriter, message string) {
	RenderError(w, http.StatusBadRequest, message)
}

func ErrorNotFound(w http.ResponseWriter, message string) {
	RenderError(w, http.StatusNotFound, message)
}

func ErrorInternalServer(w http.ResponseWriter) {
	RenderError(w, http.StatusInternalServerError, "Произошла внутренняя ошибка сервера")
}

func ErrorForbidden(w http.ResponseWriter, message string) {
	RenderError(w, http.StatusForbidden, message)
}

func ErrorUnauthorized(w http.ResponseWriter) {
	RenderError(w, http.StatusUnauthorized, "Unauthorized access. Please log in.")
}
