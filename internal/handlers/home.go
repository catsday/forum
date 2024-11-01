package handlers

import (
	"database/sql"
	"forum/internal/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type TemplateData struct {
	Posts            []*models.Post
	Username         string
	LoggedIn         bool
	ActiveCategoryID int
}

func Home(w http.ResponseWriter, r *http.Request, postModel *models.PostModel, commentModel *models.CommentModel, db *sql.DB) {
	userID, err := GetSessionUserID(r, db)
	loggedIn := err == nil

	var username string
	if loggedIn {
		err = db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	var posts []*models.Post
	activeCategoryID := 0

	if r.URL.Query().Get("likedPosts") == "1" && loggedIn {
		posts, err = postModel.GetLikedPostsByUserID(userID)
		if err != nil {
			log.Printf("Error retrieving liked posts: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else if r.URL.Query().Get("myPosts") == "1" && loggedIn {
		posts, err = postModel.GetByUserID(userID)
		if err != nil {
			log.Printf("Error retrieving user's posts: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		categoryIDStr := r.URL.Query().Get("categoryID")
		if categoryIDStr != "" {
			categoryID, convErr := strconv.Atoi(categoryIDStr)
			if convErr == nil {
				posts, err = postModel.GetByCategoryID(categoryID, userID)
				activeCategoryID = categoryID
				if err != nil {
					log.Printf("Error retrieving posts by category: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			} else {
				log.Printf("Error converting categoryID to integer: %v", convErr)
				http.Error(w, "Invalid category ID", http.StatusBadRequest)
				return
			}
		} else {
			posts, err = postModel.Latest(userID)
			if err != nil {
				log.Printf("Error retrieving latest posts: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	}

	for _, post := range posts {
		if loggedIn {
			post.UserCommented, err = commentModel.HasUserCommented(post.ID, userID)
			if err != nil {
				log.Printf("Error checking user comments: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		post.CommentCount, err = commentModel.CountByPostID(post.ID)
		if err != nil {
			log.Printf("Error counting comments: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	funcMap := template.FuncMap{
		"split": func(input string) []string {
			return strings.Split(input, "\n")
		},
	}

	files := []string{
		"./ui/templates/home.html",
	}

	ts, err := template.New("home.html").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		log.Printf("Error parsing template files: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		Posts:            posts,
		Username:         username,
		LoggedIn:         loggedIn,
		ActiveCategoryID: activeCategoryID,
	}

	if err := ts.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
