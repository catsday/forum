package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID       int
	PostID   int
	UserID   int
	Created  time.Time
	Content  string
	Username string
}

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) Insert(postID, userID int, content string) error {
	stmt := `INSERT INTO comments (post_id, user_id, content, created) VALUES (?, ?, ?, ?)`
	_, err := m.DB.Exec(stmt, postID, userID, content, time.Now().In(gmtPlus5))
	return err
}

func (m *CommentModel) GetByPostID(postID int) ([]*Comment, error) {
	stmt := `SELECT c.id, c.post_id, c.user_id, c.created, c.content, u.username
             FROM comments c
             JOIN users u ON c.user_id = u.id
             WHERE c.post_id = ? ORDER BY c.created ASC`
	rows, err := m.DB.Query(stmt, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		c := &Comment{}
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Created, &c.Content, &c.Username)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (m *CommentModel) CountByPostID(postID int) (int, error) {
	var count int
	err := m.DB.QueryRow(`SELECT COUNT(*) FROM comments WHERE post_id = ?`, postID).Scan(&count)
	return count, err
}

func (m *CommentModel) HasUserCommented(postID, userID int) (bool, error) {
	var exists bool
	err := m.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM comments WHERE post_id = ? AND user_id = ?)`, postID, userID).Scan(&exists)
	return exists, err
}