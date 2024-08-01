package internal

import (
	"forum/internal/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type templateData struct {
	Posts []*models.Post
}

func Home(w http.ResponseWriter, r *http.Request, postModel *models.PostModel) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	posts, err := postModel.Latest()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error retrieving latest posts: %v", err)
		return
	}

	files := []string{
		"./ui/templates/home.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing template files: %v", err)
		return
	}

	data := templateData{
		Posts: posts,
	}

	err = ts.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}

func PostView(w http.ResponseWriter, r *http.Request, postModel *models.PostModel) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	post, err := postModel.Get(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/templates/view.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, post)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func PostCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
