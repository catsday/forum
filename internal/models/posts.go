package models

import (
	"database/sql"
	"errors"
	"time"
)

type Post struct {
	ID      int
	Title   string
	Content string
	Created time.Time
}

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO posts (title, content, created)
    VALUES(?, ?, datetime('now', 'utc'))`

	result, err := m.DB.Exec(stmt, title, content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PostModel) Get(id int) (*Post, error) {
	stmt := `SELECT id, title, content, created FROM posts WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &Post{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created)
	if err == sql.ErrNoRows {
		return nil, errors.New("no matching record found")
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *PostModel) Latest() ([]*Post, error) {
	stmt := `SELECT id, title, content, created FROM posts ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		s := &Post{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
