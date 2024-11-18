package handlers

import (
	"database/sql"
	"forum/internal/models"
	"html/template"
	"net/http"
	"strconv"
)

type TemplateData struct {
	Posts            []*models.Post
	Username         string
	LoggedIn         bool
	ActiveCategoryID int
	FilterMyPosts    bool
	FilterLikedPosts bool
	FilterComments   bool
	SortOrder        string
	SortBy           string
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

	filterMyPosts := r.URL.Query().Get("myPosts") == "1" && loggedIn
	filterLikedPosts := r.URL.Query().Get("likedPosts") == "1" && loggedIn
	filterComments := r.URL.Query().Get("commentedPosts") == "1" && loggedIn

	sortOrder := r.URL.Query().Get("sort")
	sortBy := r.URL.Query().Get("sortBy")

	var posts []*models.Post
	activeCategoryID := 0

	if filterComments {
		posts, err = postModel.GetPostsWithUserComments(userID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else if filterLikedPosts {
		posts, err = postModel.GetLikedPostsByUserID(userID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else if filterMyPosts {
		posts, err = postModel.GetByUserID(userID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		categoryIDStr := r.URL.Query().Get("categoryID")
		if categoryIDStr != "" {
			categoryID, convErr := strconv.Atoi(categoryIDStr)
			if convErr == nil {
				if sortOrder == "asc" {
					posts, err = postModel.GetByCategoryIDAsc(categoryID, userID)
				} else {
					posts, err = postModel.GetByCategoryID(categoryID, userID)
				}
				activeCategoryID = categoryID
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Invalid category ID", http.StatusBadRequest)
				return
			}
		} else {
			if sortBy == "likes" {
				posts, err = postModel.MostLiked(userID)
			} else if sortBy == "comments" {
				posts, err = postModel.MostCommented(userID)
			} else if sortOrder == "asc" {
				posts, err = postModel.Oldest(userID)
			} else {
				posts, err = postModel.Latest(userID)
			}
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	}

	for _, post := range posts {
		if loggedIn {
			post.UserCommented, err = commentModel.HasUserCommented(post.ID, userID)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		post.CommentCount, err = commentModel.CountByPostID(post.ID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	files := []string{
		"./ui/templates/home.html",
		"./ui/templates/header.html",
		"./ui/templates/footer.html",
		"./ui/templates/left_sidebar.html",
		"./ui/templates/right_sidebar.html",
	}

	ts, err := template.New("home.html").ParseFiles(files...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		Posts:            posts,
		Username:         username,
		LoggedIn:         loggedIn,
		ActiveCategoryID: activeCategoryID,
		FilterMyPosts:    filterMyPosts,
		FilterLikedPosts: filterLikedPosts,
		FilterComments:   filterComments,
		SortOrder:        sortOrder,
		SortBy:           sortBy,
	}

	if err := ts.Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
