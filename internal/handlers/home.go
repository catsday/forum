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
}

func Home(w http.ResponseWriter, r *http.Request, postModel *models.PostModel, commentModel *models.CommentModel, db *sql.DB) {
	userModel := &models.UserModel{DB: db}
	userID, err := userModel.GetSessionUserIDFromRequest(r)
	loggedIn := err == nil

	var username string
	if loggedIn {
		err = db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Failed to connect to the database. Please try again later.")
			return
		}
	}

	if r.Method != http.MethodGet {
		RenderError(w, http.StatusMethodNotAllowed, "The requested resource does not support the HTTP method used. Please verify your request.")
		return
	}

	filterMyPosts := r.URL.Query().Get("myPosts") == "1" && loggedIn
	filterLikedPosts := r.URL.Query().Get("likedPosts") == "1" && loggedIn
	filterComments := r.URL.Query().Get("commentedPosts") == "1" && loggedIn

	var posts []*models.Post
	activeCategoryID := 0

	defer func() {
		if r := recover(); r != nil {
			RenderError(w, http.StatusInternalServerError, "The server encountered an unexpected error. Please try again later.")
		}
	}()

	if filterComments {
		posts, err = postModel.GetPostsWithUserComments(userID)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Failed to connect to the database. Please try again later.")
			return
		}
	} else if filterLikedPosts {
		posts, err = postModel.GetLikedPostsByUserID(userID)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Failed to connect to the database. Please try again later.")
			return
		}
	} else if filterMyPosts {
		posts, err = postModel.GetByUserID(userID)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Failed to connect to the database. Please try again later.")
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
					RenderError(w, http.StatusInternalServerError, "Failed to connect to the database. Please try again later.")
					return
				}
			} else {
				RenderError(w, http.StatusBadRequest, "The category ID provided is invalid. Please check your input.")
				return
			}
		} else {
			posts, err = postModel.Latest(userID)
			if err != nil {
				RenderError(w, http.StatusInternalServerError, "Failed to connect to the database. Please try again later.")
				return
			}
		}
	}

	for _, post := range posts {
		if loggedIn {
			post.UserCommented, err = commentModel.HasUserCommented(post.ID, userID)
			if err != nil {
				RenderError(w, http.StatusInternalServerError, "Failed to connect to the database. Please try again later.")
				return
			}
		}

		post.CommentCount, err = commentModel.CountByPostID(post.ID)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Failed to connect to the database. Please try again later.")
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
		RenderError(w, http.StatusInternalServerError, "The server failed to load the required template files. Please try again later.")
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
	}

	if err := ts.Execute(w, data); err != nil {
		http.Error(w, "Failed to render the page. The server encountered a technical issue.", http.StatusInternalServerError)
	}
}
