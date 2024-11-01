package handlers

import (
	"database/sql"
	"forum/internal/models"
	"net/http"
	"strconv"
)

func AddComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userModel := &models.UserModel{DB: db}
	userID, err := userModel.GetSessionUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Path[len("/post/") : len(r.URL.Path)-len("/comment")]
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	commentModel := &models.CommentModel{DB: db}
	if err := commentModel.Insert(postID, userID, content); err != nil {
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+idStr, http.StatusSeeOther)
}
