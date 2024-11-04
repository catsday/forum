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

func PostView(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	idStr := strings.TrimPrefix(r.URL.Path, "/post/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		log.Printf("Invalid post ID: %v", err)
		http.NotFound(w, r)
		return
	}

	postModel := &models.PostModel{DB: db}
	commentModel := &models.CommentModel{DB: db}

	post, err := postModel.Get(id)
	if err == sql.ErrNoRows {
		log.Printf("Post with ID %d not found", id)
		http.NotFound(w, r)
		return
	} else if err != nil {
		log.Printf("Error retrieving post with ID %d: %v", id, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("SELECT username FROM users WHERE id = ?", post.UserID).Scan(&post.Username)
	if err != nil {
		log.Printf("Error retrieving username for post ID %d: %v", post.ID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	comments, err := commentModel.GetByPostID(id)
	if err != nil {
		log.Printf("Error retrieving comments for post ID %d: %v", id, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userModel := &models.UserModel{DB: db}
	userID, _ := userModel.GetSessionUserIDFromRequest(r)
	if userID > 0 {
		post.UserVote, _ = postModel.GetUserVote(post.ID, userID)
	}

	post.Likes, post.Dislikes, err = postModel.GetLikesAndDislikes(post.ID)
	if err != nil {
		log.Printf("Error retrieving likes/dislikes for post ID %d: %v", post.ID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Post             *models.Post
		Comments         []*models.Comment
		LoggedIn         bool
		Username         string
		ActiveCategoryID int
		FilterMyPosts    bool
		FilterLikedPosts bool
		FilterComments   bool
	}{
		Post:             post,
		Comments:         comments,
		LoggedIn:         userID > 0,
		Username:         post.Username,
		ActiveCategoryID: 0,
		FilterMyPosts:    false,
		FilterLikedPosts: false,
		FilterComments:   false,
	}

	funcMap := template.FuncMap{
		"split": func(input string) []string {
			return strings.Split(input, "\n")
		},
	}

	files := []string{
		"./ui/templates/view.html",
		"./ui/templates/left_sidebar.html",
		"./ui/templates/right_sidebar.html",
	}

	ts, err := template.New("view.html").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		log.Printf("Error parsing template files: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template for post ID %d: %v", post.ID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func PostCreateForm(w http.ResponseWriter, r *http.Request) {
	files := []string{"./ui/templates/create.html"}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	ts.Execute(w, nil)
}

func PostCreate(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userModel := &models.UserModel{DB: db}
	userID, err := userModel.GetSessionUserIDFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/forum/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	if title == "" || content == "" {
		http.Error(w, "Title and Content cannot be empty", http.StatusBadRequest)
		return
	}

	var categoryIDs []int
	for _, categoryIDStr := range r.Form["categories"] {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err == nil {
			categoryIDs = append(categoryIDs, categoryID)
		}
	}

	postModel := &models.PostModel{DB: db}
	postID, err := postModel.InsertWithUserIDAndCategories(title, content, userID, categoryIDs)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
}
