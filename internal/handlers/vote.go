package handlers

import (
	"database/sql"
	"forum/internal/models"
	"net/http"
	"strconv"
)

func ToggleVote(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userModel := &models.UserModel{DB: db}
	userID, err := userModel.GetSessionUserIDFromRequest(r)
	if err != nil || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("postID"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	voteType, err := strconv.Atoi(r.FormValue("voteType"))
	if err != nil || (voteType != 1 && voteType != -1) {
		http.Error(w, "Invalid vote type", http.StatusBadRequest)
		return
	}

	postModel := &models.PostModel{DB: db}
	err = postModel.ToggleVote(postID, userID, voteType)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
