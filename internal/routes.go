package internal

import (
	"database/sql"
	"forum/internal/models"
	"net/http"
)

func Router(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	postModel := &models.PostModel{DB: db}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r, postModel)
	})

	mux.HandleFunc("/forum/login", Login(db))

	mux.HandleFunc("/forum/signup", func(w http.ResponseWriter, r *http.Request) {
		SignUp(w, r, db)
	})

	mux.HandleFunc("/forum/logout", Logout)

	mux.HandleFunc("/forum/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			PostCreateForm(w, r)
		} else {
			PostCreate(w, r, postModel)
		}
	})

	mux.HandleFunc("/forum/view", func(w http.ResponseWriter, r *http.Request) {
		PostView(w, r, postModel)
	})

	mux.HandleFunc("/forum/profile", func(w http.ResponseWriter, r *http.Request) {
		UserProfile(w, r, db)
	})

	return mux
}
