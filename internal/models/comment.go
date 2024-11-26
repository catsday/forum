package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID       int       `json:"id"`
	PostID   int       `json:"post_id"`
	UserID   int       `json:"user_id"`
	Created  time.Time `json:"created"`
	Content  string    `json:"content"`
	Username string    `json:"username"`
	Likes    int       `json:"likes"`
	Dislikes int       `json:"dislikes"`
	UserVote int       `json:"user_vote"`
}

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) Insert(postID, userID int, content string) error {
	stmt := `INSERT INTO comments (post_id, user_id, content, created) VALUES (?, ?, ?, ?)`
	_, err := m.DB.Exec(stmt, postID, userID, content, time.Now().In(gmtPlus5))
	return err
}

func (m *CommentModel) GetByPostID(postID int, userID int) ([]*Comment, error) {
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

		c.Likes, c.Dislikes, err = m.GetLikesAndDislikes(c.ID)
		if err != nil {
			return nil, err
		}

		if userID > 0 {
			c.UserVote, _ = m.GetUserVote(c.ID, userID)
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

func (m *CommentModel) ToggleVote(commentID, userID, voteType int) error {
	var existingVote int
	err := m.DB.QueryRow("SELECT vote_type FROM comment_votes WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&existingVote)
	if err == nil && existingVote == voteType {
		_, err = m.DB.Exec("DELETE FROM comment_votes WHERE comment_id = ? AND user_id = ?", commentID, userID)
		return err
	}
	if err == nil && existingVote != voteType {
		_, err = m.DB.Exec("UPDATE comment_votes SET vote_type = ? WHERE comment_id = ? AND user_id = ?", voteType, commentID, userID)
		return err
	}
	_, err = m.DB.Exec("INSERT INTO comment_votes (comment_id, user_id, vote_type) VALUES (?, ?, ?)", commentID, userID, voteType)
	return err
}

func (m *CommentModel) GetLikesAndDislikes(commentID int) (int, int, error) {
	var likes, dislikes int
	err := m.DB.QueryRow("SELECT COUNT(*) FROM comment_votes WHERE comment_id = ? AND vote_type = 1", commentID).Scan(&likes)
	if err != nil {
		return 0, 0, err
	}
	err = m.DB.QueryRow("SELECT COUNT(*) FROM comment_votes WHERE comment_id = ? AND vote_type = -1", commentID).Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}
	return likes, dislikes, nil
}

func (m *CommentModel) GetUserVote(commentID, userID int) (int, error) {
	var voteType int
	err := m.DB.QueryRow("SELECT vote_type FROM comment_votes WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&voteType)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return voteType, nil
}
