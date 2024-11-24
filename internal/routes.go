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

	mux.HandleFunc("/forum/technology", func(w http.ResponseWriter, r *http.Request) {
		r.URL.RawQuery = "categoryID=1"
		handlers.Home(w, r, postModel, commentModel, db)
	})
	mux.HandleFunc("/forum/entertainment", func(w http.ResponseWriter, r *http.Request) {
		r.URL.RawQuery = "categoryID=2"
		handlers.Home(w, r, postModel, commentModel, db)
	})
	mux.HandleFunc("/forum/sports", func(w http.ResponseWriter, r *http.Request) {
		r.URL.RawQuery = "categoryID=3"
		handlers.Home(w, r, postModel, commentModel, db)
	})
	mux.HandleFunc("/forum/education", func(w http.ResponseWriter, r *http.Request) {
		r.URL.RawQuery = "categoryID=4"
		handlers.Home(w, r, postModel, commentModel, db)
	})
	mux.HandleFunc("/forum/health", func(w http.ResponseWriter, r *http.Request) {
		r.URL.RawQuery = "categoryID=5"
		handlers.Home(w, r, postModel, commentModel, db)
	})

	mux.HandleFunc("/forum/posted", handlers.AuthorizeAndHandle(db, func(w http.ResponseWriter, r *http.Request, userID int) {
		r.URL.RawQuery = "myPosts=1"
		handlers.Home(w, r, postModel, commentModel, db)
	}))

	mux.HandleFunc("/forum/liked", handlers.AuthorizeAndHandle(db, func(w http.ResponseWriter, r *http.Request, userID int) {
		r.URL.RawQuery = "likedPosts=1"
		handlers.Home(w, r, postModel, commentModel, db)
	}))

	mux.HandleFunc("/forum/commented", handlers.AuthorizeAndHandle(db, func(w http.ResponseWriter, r *http.Request, userID int) {
		r.URL.RawQuery = "commentedPosts=1"
		handlers.Home(w, r, postModel, commentModel, db)
	}))

	mux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/comment") {
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
			handlers.PostCreateForm(w, r, db)
		} else {
			handlers.PostCreate(w, r, db)
		}
	})
	mux.HandleFunc("/toggle-vote", func(w http.ResponseWriter, r *http.Request) {
		handlers.ToggleVote(w, r, db)
	})

	mux.HandleFunc("/forum/toggle-ban", func(w http.ResponseWriter, r *http.Request) {
		handlers.ToggleBanStatus(w, r, db)
	})
	mux.HandleFunc("/forum/profile/change-password", func(w http.ResponseWriter, r *http.Request) {
		handlers.ChangePassword(db).ServeHTTP(w, r)
	})
	mux.HandleFunc("/forum/profile/change-name", func(w http.ResponseWriter, r *http.Request) {
		handlers.ChangeName(db).ServeHTTP(w, r)
	})

	return mux
}
