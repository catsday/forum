package internal

import (
	"database/sql"
	"forum/internal/handlers"
	"forum/internal/models"
	"net/http"
	"strings"
)

func Router(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	postModel := &models.PostModel{DB: db}
	commentModel := &models.CommentModel{DB: db}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.Home(w, r, postModel, commentModel, db)
	})

	mux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/post/") && strings.HasSuffix(r.URL.Path, "/comment") {
			handlers.AddComment(w, r, db)
		} else {
			handlers.PostView(w, r, db)
		}
	})
	mux.HandleFunc("/forum/profile", func(w http.ResponseWriter, r *http.Request) {
		handlers.UserProfile(w, r, db)
	})
	mux.HandleFunc("/forum/login", handlers.Login(db))

	mux.HandleFunc("/forum/signup", func(w http.ResponseWriter, r *http.Request) {
		handlers.SignUp(w, r, db)
	})

	mux.HandleFunc("/forum/logout", func(w http.ResponseWriter, r *http.Request) {
		handlers.Logout(w, r, db)
	})

	mux.HandleFunc("/forum/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.PostCreateForm(w, r)
		} else {
			handlers.PostCreate(w, r, db)
		}
	})

	mux.HandleFunc("/toggle-vote", func(w http.ResponseWriter, r *http.Request) {
		handlers.ToggleVote(w, r, db)
	})

	return mux
}
