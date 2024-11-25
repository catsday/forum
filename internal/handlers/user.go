package handlers

import (
	"database/sql"
	"forum/internal/models"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"
)

type AdminUser struct {
	ID                    int
	Username              string
	Email                 string
	PostCount             int
	CommentCount          int
	LikedPosts            int
	DislikedPosts         int
	LikeDislikeRatioPosts float64
	IsBanned              bool
}

type ProfileData struct {
	ID                    int
	Username              string
	Email                 string
	PostCount             int
	CommentCount          int
	LikedPosts            int
	DislikedPosts         int
	LikeDislikeRatioPosts float64
	IsAdmin               bool
	LoggedIn              bool
	FilterMyPosts         bool
	FilterLikedPosts      bool
	FilterComments        bool
	ActiveCategoryID      int
	Users                 []AdminUser
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files := []string{
			"./ui/templates/login.html",
			"./ui/templates/header.html",
			"./ui/templates/footer.html",
			"./ui/templates/left_sidebar.html",
			"./ui/templates/right_sidebar.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Failed to load login page templates.")
			return
		}

		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			password := r.FormValue("password")

			var hashedPassword string
			var isBanned bool

			err := db.QueryRow("SELECT password, is_banned FROM users WHERE email = ?", email).Scan(&hashedPassword, &isBanned)
			if err != nil {
				if err == sql.ErrNoRows {
					RenderError(w, http.StatusUnauthorized, "The email does not exist in the database.")
					return
				}
				RenderError(w, http.StatusInternalServerError, "Failed to query user information.")
				return
			}

			if isBanned {
				RenderError(w, http.StatusForbidden, "The account is banned. Please contact support.")
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
			if err != nil {
				RenderError(w, http.StatusUnauthorized, "The entered password is incorrect.")
				return
			}

			userModel := &models.UserModel{DB: db}
			userID, err := userModel.Authenticate(email, password)
			if err != nil {
				RenderError(w, http.StatusInternalServerError, "Failed to authenticate user.")
				return
			}

			sessionID, err := userModel.CreateSession(userID)
			if err != nil {
				RenderError(w, http.StatusInternalServerError, "Failed to create user session.")
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    sessionID,
				Expires:  time.Now().Add(time.Hour),
				HttpOnly: true,
				Path:     "/",
			})

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if r.Method == http.MethodGet {
			err := ts.Execute(w, nil)
			if err != nil {
				RenderError(w, http.StatusInternalServerError, "Failed to render the login page.")
			}
			return
		}

		w.Header().Set("Allow", "GET, POST")
		RenderError(w, http.StatusMethodNotAllowed, "Method not allowed. Use GET or POST.")
	}
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}

func emailExists(db *sql.DB, email string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func SignUp(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	files := []string{
		"./ui/templates/signup.html",
		"./ui/templates/header.html",
		"./ui/templates/footer.html",
		"./ui/templates/left_sidebar.html",
		"./ui/templates/right_sidebar.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Failed to load templates for registration page.")
		return
	}

	if r.Method == http.MethodGet {
		err = ts.Execute(w, nil)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "The registration page could not be displayed.")
		}
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm-password")

		if !isValidEmail(email) {
			RenderError(w, http.StatusBadRequest, "Incorrect email. Check the entered data.")
			return
		}

		if emailExists(db, email) {
			RenderError(w, http.StatusConflict, "The email you entered is already registered. Please use another email.")
			return
		}

		if password != confirmPassword {
			RenderError(w, http.StatusBadRequest, "The 'Password' and 'Confirm password' fields must match.")
			return
		}

		if len(password) < 8 {
			RenderError(w, http.StatusBadRequest, "Password must be at least 8 characters.")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Failed to process password. Try again.")
			return
		}

		_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, string(hashedPassword))
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Failed to create user due to internal server error.")
			return
		}

		http.Redirect(w, r, "/forum/login", http.StatusSeeOther)
		return
	}

	w.Header().Set("Allow", "GET, POST")
	RenderError(w, http.StatusMethodNotAllowed, "Method not supported. Use GET or POST.")
}

func Logout(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/forum/login", http.StatusSeeOther)
		return
	}

	userModel := &models.UserModel{DB: db}
	err = userModel.DeleteSession(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	})
	http.Redirect(w, r, "/forum/login", http.StatusSeeOther)
}

func UserProfile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userModel := &models.UserModel{DB: db}

	userID, err := userModel.GetSessionUserIDFromRequest(r)
	if err != nil || userID == 0 {
		log.Printf("UserProfile: Unauthorized access attempt. Error: %v", err)
		http.Redirect(w, r, "/forum/login", http.StatusSeeOther)
		return
	}

	var username, email string
	err = db.QueryRow("SELECT username, email FROM users WHERE id = ?", userID).Scan(&username, &email)
	if err != nil {
		log.Printf("UserProfile: Failed to fetch user data for ID %d. Error: %v", userID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	isAdmin := email == "admin@gmail.com" && username == "Admin"

	var postCount, commentCount, likedPosts, dislikedPosts int
	var likeDislikeRatioPosts float64

	err = db.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&postCount)
	if err != nil {
		log.Printf("UserProfile: Failed to count posts for user ID %d. Error: %v", userID, err)
		postCount = 0
	}

	err = db.QueryRow("SELECT COUNT(*) FROM comments WHERE user_id = ?", userID).Scan(&commentCount)
	if err != nil {
		log.Printf("UserProfile: Failed to count comments for user ID %d. Error: %v", userID, err)
		commentCount = 0
	}

	err = db.QueryRow("SELECT COUNT(*) FROM post_votes WHERE user_id = ? AND vote_type = 1", userID).Scan(&likedPosts)
	if err != nil {
		log.Printf("UserProfile: Failed to count liked posts for user ID %d. Error: %v", userID, err)
		likedPosts = 0
	}

	err = db.QueryRow("SELECT COUNT(*) FROM post_votes WHERE user_id = ? AND vote_type = -1", userID).Scan(&dislikedPosts)
	if err != nil {
		log.Printf("UserProfile: Failed to count disliked posts for user ID %d. Error: %v", userID, err)
		dislikedPosts = 0
	}

	err = db.QueryRow("SELECT COALESCE(AVG(vote_type), 0) FROM post_votes WHERE user_id = ?", userID).Scan(&likeDislikeRatioPosts)
	if err != nil {
		log.Printf("UserProfile: Failed to calculate like/dislike ratio for posts for user ID %d. Error: %v", userID, err)
		likeDislikeRatioPosts = 0
	}

	var users []AdminUser
	if isAdmin {
		rows, err := db.Query(`
			SELECT 
				u.id, u.username, u.email, u.is_banned,
				(SELECT COUNT(*) FROM posts WHERE user_id = u.id) AS post_count,
				(SELECT COUNT(*) FROM comments WHERE user_id = u.id) AS comment_count,
				(SELECT COUNT(*) FROM post_votes WHERE user_id = u.id AND vote_type = 1) AS liked_posts,
				(SELECT COUNT(*) FROM post_votes WHERE user_id = u.id AND vote_type = -1) AS disliked_posts,
				COALESCE((SELECT AVG(vote_type) FROM post_votes WHERE user_id = u.id), 0) AS like_dislike_ratio_posts
			FROM users u
		`)
		if err != nil {
			log.Printf("UserProfile: Failed to fetch user list for admin. Error: %v", err)
		} else {
			defer rows.Close()
			for rows.Next() {
				var user AdminUser
				if err := rows.Scan(
					&user.ID, &user.Username, &user.Email, &user.IsBanned,
					&user.PostCount, &user.CommentCount,
					&user.LikedPosts, &user.DislikedPosts,
					&user.LikeDislikeRatioPosts,
				); err != nil {
					log.Printf("UserProfile: Error scanning user data: %v", err)
					continue
				}
				users = append(users, user)
			}
		}
	}

	data := ProfileData{
		ID:                    userID,
		Username:              username,
		Email:                 email,
		PostCount:             postCount,
		CommentCount:          commentCount,
		LikedPosts:            likedPosts,
		DislikedPosts:         dislikedPosts,
		LikeDislikeRatioPosts: likeDislikeRatioPosts,
		IsAdmin:               isAdmin,
		LoggedIn:              true,
		FilterMyPosts:         false,
		FilterLikedPosts:      false,
		FilterComments:        false,
		ActiveCategoryID:      0,
		Users:                 users,
	}

	files := []string{
		"./ui/templates/profile.html",
		"./ui/templates/header.html",
		"./ui/templates/footer.html",
		"./ui/templates/left_sidebar.html",
		"./ui/templates/right_sidebar.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("UserProfile: Failed to parse templates for user ID %d. Error: %v", userID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Printf("UserProfile: Failed to execute template for user ID %d. Error: %v", userID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("UserProfile: Successfully rendered profile page for user ID %d.", userID)
}

func GetSessionUserID(r *http.Request, db *sql.DB) (int, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, err
	}
	userModel := &models.UserModel{DB: db}
	return userModel.GetSessionUserID(cookie.Value)
}

func ToggleBanStatus(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var isBanned bool
	err := db.QueryRow("SELECT is_banned FROM users WHERE id = ?", userID).Scan(&isBanned)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("ToggleBanStatus: Error fetching user status: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	newStatus := !isBanned
	_, err = db.Exec("UPDATE users SET is_banned = ? WHERE id = ?", newStatus, userID)
	if err != nil {
		log.Printf("ToggleBanStatus: Error updating ban status: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("ToggleBanStatus: Updated ban status for user ID %s to %t", userID, newStatus)
}

func ChangePassword(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			RenderError(w, http.StatusMethodNotAllowed, "Method Not Allowed. Use POST.")
			return
		}

		userID, err := GetSessionUserID(r, db)
		if err != nil {
			RenderError(w, http.StatusUnauthorized, "Unauthorized. Please log in to change your password.")
			return
		}

		currentPassword := r.FormValue("current-password")
		newPassword := r.FormValue("new-password")
		confirmPassword := r.FormValue("confirm-password")

		var hashedPassword string
		err = db.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&hashedPassword)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Failed to retrieve your current password.")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currentPassword))
		if err != nil {
			RenderError(w, http.StatusUnauthorized, "The current password you entered is incorrect.")
			return
		}

		if len(newPassword) < 8 {
			RenderError(w, http.StatusBadRequest, "The new password must be at least 8 characters long.")
			return
		}

		if newPassword != confirmPassword {
			RenderError(w, http.StatusBadRequest, "The 'New Password' and 'Confirm Password' fields must match.")
			return
		}

		newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Failed to process the new password.")
			return
		}

		_, err = db.Exec("UPDATE users SET password = ? WHERE id = ?", newHashedPassword, userID)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Failed to update the password. Please try again.")
			return
		}

		http.Redirect(w, r, "/forum/profile", http.StatusSeeOther)
	}
}

func ChangeName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			RenderError(w, http.StatusMethodNotAllowed, "Method Not Allowed. Use POST.")
			return
		}

		userID, err := GetSessionUserID(r, db)
		if err != nil {
			RenderError(w, http.StatusUnauthorized, "Unauthorized. Please log in to change your name.")
			return
		}

		newName := r.FormValue("new-name")
		if IsBlankOrInvisible(newName) {
			RenderError(w, http.StatusBadRequest, "The name cannot consist only of invisible characters.")
			return
		}

		_, err = db.Exec("UPDATE users SET username = ? WHERE id = ?", newName, userID)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Failed to update the name. Please try again.")
			return
		}

		http.Redirect(w, r, "/forum/profile", http.StatusSeeOther)
	}
}
