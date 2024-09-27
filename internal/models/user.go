package models

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

type UserModel struct {
	DB *sql.DB
}

//func NewUserModel(db *sql.DB) *UserModel {
//	return &UserModel{DB: db}
//}

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
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}
	return id, nil
}
