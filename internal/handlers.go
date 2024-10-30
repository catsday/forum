package internal

import (
	"database/sql"
	"forum/internal/models"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

type templateData struct {
	Posts            []*models.Post
	Username         string
	LoggedIn         bool
	ActiveCategoryID int
}

func Home(w http.ResponseWriter, r *http.Request, postModel *models.PostModel) {
	session, _ := store.Get(r, "session-name")
	var username string
	var loggedIn bool
	var userID int

	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		userID = session.Values["userID"].(int)
		err := postModel.DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		loggedIn = true
	} else {
		userID = 0
	}

	var posts []*models.Post
	var err error
	activeCategoryID := 0

	if r.URL.Query().Get("likedPosts") == "1" && loggedIn {
		posts, err = postModel.GetLikedPostsByUserID(userID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error retrieving liked posts: %v", err)
			return
		}
	} else {
		categoryIDStr := r.URL.Query().Get("categoryID")
		if categoryIDStr != "" {
			categoryID, convErr := strconv.Atoi(categoryIDStr)
			if convErr == nil {
				posts, err = postModel.GetByCategoryID(categoryID)
				activeCategoryID = categoryID
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					log.Printf("Error retrieving posts by category: %v", err)
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
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Error retrieving latest posts: %v", err)
				return
			}
		}
	}

	files := []string{
		"./ui/templates/home.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing template files: %v", err)
		return
	}

	data := templateData{
		Posts:            posts,
		Username:         username,
		LoggedIn:         loggedIn,
		ActiveCategoryID: activeCategoryID,
	}

	err = ts.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}

func PostView(w http.ResponseWriter, r *http.Request, postModel *models.PostModel) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	post, err := postModel.Get(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/templates/view.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, post)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func PostCreateForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/templates/create.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing template files: %v", err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}

func PostCreate(w http.ResponseWriter, r *http.Request, postModel *models.PostModel) {
	session, _ := store.Get(r, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/forum/login", http.StatusSeeOther)
		return
	}

	userID := session.Values["userID"]

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
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

	postID, err := postModel.InsertWithUserIDAndCategories(title, content, userID.(int), categoryIDs)
	if err != nil {
		log.Printf("Error creating new post: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("New post created with ID %d and categories: %v", postID, categoryIDs)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			password := r.FormValue("password")

			var id int
			var hashedPassword string
			err := db.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&id, &hashedPassword)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "Invalid username or password", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}

			session, _ := store.Get(r, "session-name")
			session.Values["authenticated"] = true
			session.Values["userID"] = id
			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			files := []string{"./ui/templates/login.html"}
			ts, err := template.ParseFiles(files...)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			ts.Execute(w, nil)
		}
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
		log.Printf("Error checking if email exists: %v", err)
		return false
	}
	return exists
}

func SignUp(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodGet {
		files := []string{"./ui/templates/signup.html"}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}
		ts.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm-password")

		if !isValidEmail(email) {
			http.Error(w, "Invalid email address", http.StatusBadRequest)
			return
		}

		if emailExists(db, email) {
			http.Error(w, "Email already in use", http.StatusBadRequest)
			return
		}

		if len(password) < 8 {
			http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
			return
		}

		if password != confirmPassword {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, string(hashedPassword))
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/forum/login", http.StatusSeeOther)
	} else {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/forum/login", http.StatusSeeOther)
}

func UserProfile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	session, _ := store.Get(r, "session-name")

	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var username, email string
	err := db.QueryRow("SELECT username, email FROM users WHERE id = ?", userID).Scan(&username, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		ID       int
		Username string
		Email    string
	}{
		ID:       userID,
		Username: username,
		Email:    email,
	}

	files := []string{"./ui/templates/profile.html"}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	ts.Execute(w, data)
}

func init() {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
	}
}

func ToggleVote(w http.ResponseWriter, r *http.Request, postModel *models.PostModel) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	userID, ok := session.Values["userID"].(int)
	if !ok {
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

	err = postModel.ToggleVote(postID, userID, voteType)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
