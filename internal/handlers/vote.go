package handlers

import (
	"database/sql"
	"forum/internal/models"
	"net/http"
	"strconv"
)

func ToggleVote(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		RenderError(w, http.StatusMethodNotAllowed, "Method Not Allowed. Use POST.")
		return
	}

	userModel := &models.UserModel{DB: db}
	userID, err := userModel.GetSessionUserIDFromRequest(r)
	if err != nil || userID == 0 {
		RenderError(w, http.StatusUnauthorized, "Unauthorized. Please log in to vote.")
		return
	}

	postID, err := strconv.Atoi(r.FormValue("postID"))
	if err != nil || postID < 1 {
		RenderError(w, http.StatusBadRequest, "Invalid post ID.")
		return
	}

	voteType, err := strconv.Atoi(r.FormValue("voteType"))
	if err != nil || (voteType != 1 && voteType != -1) {
		RenderError(w, http.StatusBadRequest, "Vote type must be 1 (like) or -1 (dislike). Invalid value provided.")
		return
	}

	postModel := &models.PostModel{DB: db}

	_, err = postModel.Get(postID)
	if err == sql.ErrNoRows {
		RenderError(w, http.StatusNotFound, "The post with the specified ID does not exist. Please check the ID.")
		return
	} else if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to retrieve post for voting.")
		return
	}

	err = postModel.ToggleVote(postID, userID, voteType)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to process your vote.")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ToggleCommentVote(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		RenderError(w, http.StatusMethodNotAllowed, "Method Not Allowed. Use POST.")
		return
	}

	userModel := &models.UserModel{DB: db}
	userID, err := userModel.GetSessionUserIDFromRequest(r)
	if err != nil || userID == 0 {
		RenderError(w, http.StatusUnauthorized, "Unauthorized. Please log in to vote.")
		return
	}

	commentID, err := strconv.Atoi(r.FormValue("commentID"))
	if err != nil || commentID < 1 {
		RenderError(w, http.StatusBadRequest, "Invalid comment ID.")
		return
	}

	voteType, err := strconv.Atoi(r.FormValue("voteType"))
	if err != nil || (voteType != 1 && voteType != -1) {
		RenderError(w, http.StatusBadRequest, "Vote type must be 1 (like) or -1 (dislike). Invalid value provided.")
		return
	}

	commentModel := &models.CommentModel{DB: db}

	_, err = db.Exec("SELECT id FROM comments WHERE id = ?", commentID)
	if err == sql.ErrNoRows {
		RenderError(w, http.StatusNotFound, "The comment with the specified ID does not exist.")
		return
	} else if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to retrieve comment for voting.")
		return
	}

	err = commentModel.ToggleVote(commentID, userID, voteType)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to process your vote.")
		return
	}

	w.WriteHeader(http.StatusOK)
}
