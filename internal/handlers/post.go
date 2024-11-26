package handlers

import (
	"database/sql"
	"forum/internal/models"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func PostView(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	idStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/post/"))

	if idStr == "" {
		RenderError(w, http.StatusBadRequest, "Id cannot be empty. Input is required.")
		return
	}

	if strings.Contains(idStr, " ") || idStr != strings.TrimSpace(idStr) {
		RenderError(w, http.StatusBadRequest, "Id cannot contain spaces. Input must be an integer.")
		return
	}

	if !regexp.MustCompile(`^\d+$`).MatchString(idStr) {
		RenderError(w, http.StatusBadRequest, "Id contains invalid characters. Only numerical values are allowed.")
		return
	}

	if len(idStr) > 18 {
		RenderError(w, http.StatusBadRequest, "Id is too long. Input must be shorter.")
		return
	}

	if len(idStr) > 1 && idStr[0] == '0' {
		RenderError(w, http.StatusBadRequest, "Id has an invalid format. While numerically equivalent to 1, the formatting is incorrect.")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		RenderError(w, http.StatusBadRequest, "Id must be a number.")
		return
	}

	if id < 1 {
		if id == 0 {
			RenderError(w, http.StatusBadRequest, "Id cannot be zero.")
		} else {
			RenderError(w, http.StatusBadRequest, "Negative Id values are not allowed. This is a client error.")
		}
		return
	}

	const MaxID = 1_000_000_000
	if id > MaxID {
		RenderError(w, http.StatusNotFound, "Id is out of the allowable range.")
		return
	}

	postModel := &models.PostModel{DB: db}

	post, err := postModel.Get(id)
	if err == sql.ErrNoRows {
		RenderError(w, http.StatusNotFound, "The post with the specified ID does not exist. Please check the ID.")
		return
	} else if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to retrieve the post.")
		return
	}

	commentModel := &models.CommentModel{DB: db}

	err = db.QueryRow("SELECT username FROM users WHERE id = ?", post.UserID).Scan(&post.Username)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to retrieve the post author's username.")
		return
	}

	userModel := &models.UserModel{DB: db}
	userID, _ := userModel.GetSessionUserIDFromRequest(r)

	comments, err := commentModel.GetByPostID(id, userID)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to retrieve comments for the post.")
		return
	}

	post.Likes, post.Dislikes, err = postModel.GetLikesAndDislikes(post.ID)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to retrieve likes and dislikes for the post.")
		return
	}

	if userID > 0 {
		post.UserVote, _ = postModel.GetUserVote(post.ID, userID)
	}

	rows, err := db.Query(`
        SELECT c.name 
        FROM categories c
        JOIN post_categories pc ON c.id = pc.category_id
        WHERE pc.post_id = ?`, post.ID)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to retrieve categories for the post.")
		return
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			RenderError(w, http.StatusInternalServerError, "Failed to retrieve category for the post.")
			return
		}
		categories = append(categories, category)
	}

	post.Categories = categories

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

	files := []string{
		"./ui/templates/view.html",
		"./ui/templates/header.html",
		"./ui/templates/footer.html",
		"./ui/templates/left_sidebar.html",
		"./ui/templates/right_sidebar.html",
	}

	ts, err := template.New("view.html").ParseFiles(files...)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to load templates for the post view.")
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to render the post view.")
	}
}

func PostCreateForm(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userModel := &models.UserModel{DB: db}
	userID, err := userModel.GetSessionUserIDFromRequest(r)
	if err != nil {
		RenderError(w, http.StatusUnauthorized, "Only authorized users can create posts. Please log in.")
		return
	}

	var username string
	err = db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to retrieve username for the session.")
		return
	}

	data := struct {
		LoggedIn         bool
		Username         string
		FilterMyPosts    bool
		FilterLikedPosts bool
		FilterComments   bool
		ActiveCategoryID int
	}{
		LoggedIn:         true,
		Username:         username,
		FilterMyPosts:    false,
		FilterLikedPosts: false,
		FilterComments:   false,
		ActiveCategoryID: 0,
	}

	files := []string{
		"./ui/templates/create.html",
		"./ui/templates/header.html",
		"./ui/templates/footer.html",
		"./ui/templates/left_sidebar.html",
		"./ui/templates/right_sidebar.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to load templates for the post creation form.")
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to render the post creation form.")
	}
}

func PostCreate(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userModel := &models.UserModel{DB: db}
	userID, err := userModel.GetSessionUserIDFromRequest(r)
	if err != nil {
		RenderError(w, http.StatusUnauthorized, "Only authorized users can create posts. Please log in.")
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		RenderError(w, http.StatusMethodNotAllowed, "Method Not Allowed. Use POST.")
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	if IsBlankOrInvisible(title) || IsBlankOrInvisible(content) {
		RenderError(w, http.StatusBadRequest, "Title and content cannot consist only of invisible characters.")
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
		RenderError(w, http.StatusInternalServerError, "Failed to create the post due to an internal error.")
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
}
