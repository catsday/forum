package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Create(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	_, err = m.DB.Exec(stmt, username, email, string(hashedPassword))
	return err
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword string
	stmt := `SELECT id, password FROM users WHERE email = ?`
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("invalid credentials")
		}
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return 0, errors.New("invalid credentials")
		}
		return 0, err
	}
	return id, nil
}

func (m *UserModel) CreateSession(userID int) (string, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return "", err
	}
	expiry := time.Now().Add(24 * time.Hour)

	stmt := `REPLACE INTO sessions (session_id, user_id, expiry) VALUES (?, ?, ?)`
	_, err = m.DB.Exec(stmt, sessionID, userID, expiry)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (m *UserModel) DeleteSession(sessionID string) error {
	stmt := `DELETE FROM sessions WHERE session_id = ?`
	_, err := m.DB.Exec(stmt, sessionID)
	return err
}

func (m *UserModel) GetSessionUserID(sessionID string) (int, error) {
	var userID int
	var expiry time.Time
	stmt := `SELECT user_id, expiry FROM sessions WHERE session_id = ?`
	err := m.DB.QueryRow(stmt, sessionID).Scan(&userID, &expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("session not found")
		}
		return 0, err
	}
	if time.Now().After(expiry) {
		_ = m.DeleteSession(sessionID)
		return 0, errors.New("session expired")
	}
	return userID, nil
}

func generateSessionID() (string, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}

func (m *UserModel) GetSessionUserIDFromRequest(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, err
	}
	return m.GetSessionUserID(cookie.Value)
}
