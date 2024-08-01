package internal

import (
	"database/sql"
	"net/http"
	"forum/internal/models"
)

func Router(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	postModel := &models.PostModel{DB: db}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r, postModel)
	})
	mux.HandleFunc("/forum/view", func(w http.ResponseWriter, r *http.Request) {
		PostView(w, r, postModel)
	})
	mux.HandleFunc("/forum/create", PostCreate)

	return mux
}
