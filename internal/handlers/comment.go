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
		RenderError(w, http.StatusUnauthorized, "Unauthorized. Please log in to add a comment.")
		return
	}

	idStr := r.URL.Path[len("/post/") : len(r.URL.Path)-len("/comment")]
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		RenderError(w, http.StatusBadRequest, "Invalid post ID. The ID must be a valid number.")
		return
	}

	content := r.FormValue("content")
	if IsBlankOrInvisible(content) {
		RenderError(w, http.StatusBadRequest, "Content cannot consist only of invisible characters.")
		return
	}

	commentModel := &models.CommentModel{DB: db}
	if err := commentModel.Insert(postID, userID, content); err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to add the comment due to an internal error.")
		return
	}

	http.Redirect(w, r, "/post/"+idStr, http.StatusSeeOther)
}
