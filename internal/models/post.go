package models

import (
	"database/sql"
	"errors"
	"time"
)

type Post struct {
	ID            int
	Title         string
	Content       string
	Created       time.Time
	UserID        int
	Username      string
	Likes         int
	Dislikes      int
	UserVote      int
	Categories    []string
	UserCommented bool
	CommentCount  int
	UserComments  []*Comment
}

type PostModel struct {
	DB *sql.DB
}

var gmtPlus5 = time.FixedZone("GMT+5", 5*60*60)

func (m *PostModel) InsertWithUserIDAndCategories(title string, content string, userID int, categoryIDs []int) (int, error) {
	if len(title) > 25 {
		title = title[:25]
	}
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO posts (title, content, user_id, created) VALUES (?, ?, ?, ?)`
	result, err := tx.Exec(stmt, title, content, userID, time.Now().In(gmtPlus5))
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, categoryID := range categoryIDs {
		_, err := tx.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, categoryID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

func (m *PostModel) Get(id int) (*Post, error) {
	query := `SELECT id, title, content, user_id, created FROM posts WHERE id = ?`
	row := m.DB.QueryRow(query, id)

	post := &Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return post, nil
}

func (m *PostModel) GetCategories(postID int) ([]string, error) {
	stmt := `SELECT categories.name FROM categories
             JOIN post_categories ON categories.id = post_categories.category_id
             WHERE post_categories.post_id = ?`

	rows, err := m.DB.Query(stmt, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (m *PostModel) Latest(userID int) ([]*Post, error) {
	stmt := `
        SELECT posts.id, posts.title, posts.content, posts.created, users.username
        FROM posts
        JOIN users ON posts.user_id = users.id
        ORDER BY posts.created DESC LIMIT 10
    `

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created, &post.Username)
		if err != nil {
			return nil, err
		}

		post.Categories, err = m.GetCategories(post.ID)
		if err != nil {
			return nil, err
		}

		post.Likes, post.Dislikes, err = m.GetLikesAndDislikes(post.ID)
		if err != nil {
			return nil, err
		}

		err = m.DB.QueryRow("SELECT vote_type FROM post_votes WHERE post_id = ? AND user_id = ?", post.ID, userID).Scan(&post.UserVote)
		if err == sql.ErrNoRows {
			post.UserVote = 0
		} else if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) GetByCategoryID(categoryID, userID int) ([]*Post, error) {
	stmt := `
        SELECT posts.id, posts.title, posts.content, posts.created, users.username
        FROM posts
        JOIN users ON posts.user_id = users.id
        JOIN post_categories ON posts.id = post_categories.post_id
        WHERE post_categories.category_id = ?
        ORDER BY posts.created DESC
    `

	rows, err := m.DB.Query(stmt, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created, &post.Username)
		if err != nil {
			return nil, err
		}

		post.Categories, err = m.GetCategories(post.ID)
		if err != nil {
			return nil, err
		}

		post.Likes, post.Dislikes, err = m.GetLikesAndDislikes(post.ID)
		if err != nil {
			return nil, err
		}

		err = m.DB.QueryRow("SELECT vote_type FROM post_votes WHERE post_id = ? AND user_id = ?", post.ID, userID).Scan(&post.UserVote)
		if err == sql.ErrNoRows {
			post.UserVote = 0
		} else if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) GetByUserID(userID int) ([]*Post, error) {
	stmt := `
        SELECT posts.id, posts.title, posts.content, posts.created, users.username
        FROM posts
        JOIN users ON posts.user_id = users.id
        WHERE posts.user_id = ?
        ORDER BY posts.created DESC
    `

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created, &post.Username)
		if err != nil {
			return nil, err
		}

		post.Categories, err = m.GetCategories(post.ID)
		if err != nil {
			return nil, err
		}

		post.Likes, post.Dislikes, err = m.GetLikesAndDislikes(post.ID)
		if err != nil {
			return nil, err
		}

		err = m.DB.QueryRow("SELECT vote_type FROM post_votes WHERE post_id = ? AND user_id = ?", post.ID, userID).Scan(&post.UserVote)
		if err == sql.ErrNoRows {
			post.UserVote = 0
		} else if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) ToggleVote(postID, userID, voteType int) error {
	var existingVote int
	err := m.DB.QueryRow("SELECT vote_type FROM post_votes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&existingVote)
	if err == nil && existingVote == voteType {
		_, err = m.DB.Exec("DELETE FROM post_votes WHERE post_id = ? AND user_id = ?", postID, userID)
		return err
	}
	if err == nil && existingVote != voteType {
		_, err = m.DB.Exec("UPDATE post_votes SET vote_type = ? WHERE post_id = ? AND user_id = ?", voteType, postID, userID)
		return err
	}
	_, err = m.DB.Exec("INSERT INTO post_votes (post_id, user_id, vote_type) VALUES (?, ?, ?)", postID, userID, voteType)
	return err
}

func (m *PostModel) GetLikesAndDislikes(postID int) (int, int, error) {
	var likes, dislikes int
	err := m.DB.QueryRow("SELECT COUNT(*) FROM post_votes WHERE post_id = ? AND vote_type = 1", postID).Scan(&likes)
	if err != nil {
		return 0, 0, err
	}
	err = m.DB.QueryRow("SELECT COUNT(*) FROM post_votes WHERE post_id = ? AND vote_type = -1", postID).Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}
	return likes, dislikes, nil
}

func (m *PostModel) GetLikedPostsByUserID(userID int) ([]*Post, error) {
	stmt := `
        SELECT posts.id, posts.title, posts.content, posts.created, users.username
        FROM posts
        JOIN users ON posts.user_id = users.id
        JOIN post_votes ON posts.id = post_votes.post_id
        WHERE post_votes.user_id = ? AND post_votes.vote_type = 1
        ORDER BY posts.created DESC
    `

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created, &post.Username)
		if err != nil {
			return nil, err
		}
		post.UserVote = 1
		post.Categories, err = m.GetCategories(post.ID)
		if err != nil {
			return nil, err
		}
		post.Likes, post.Dislikes, err = m.GetLikesAndDislikes(post.ID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) GetUserVote(postID, userID int) (int, error) {
	var voteType int
	err := m.DB.QueryRow(`SELECT vote_type FROM post_votes WHERE post_id = ? AND user_id = ?`, postID, userID).Scan(&voteType)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return voteType, nil
}

func (m *PostModel) GetUsername(userID int) (string, error) {
	var username string
	stmt := `SELECT username FROM users WHERE id = ?`
	err := m.DB.QueryRow(stmt, userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", err
	}
	return username, nil
}

func (m *PostModel) GetPostsWithUserComments(userID int) ([]*Post, error) {
	stmt := `
		SELECT DISTINCT posts.id, posts.title, posts.content, posts.created, users.username
		FROM posts
		JOIN comments ON posts.id = comments.post_id
		JOIN users ON posts.user_id = users.id
		WHERE comments.user_id = ?
		ORDER BY posts.created DESC
	`
	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created, &post.Username)
		if err != nil {
			return nil, err
		}

		post.UserComments, err = m.GetCommentsByUserIDForPost(post.ID, userID)
		if err != nil {
			return nil, err
		}

		post.Likes, post.Dislikes, err = m.GetLikesAndDislikes(post.ID)
		if err != nil {
			return nil, err
		}

		err = m.DB.QueryRow("SELECT vote_type FROM post_votes WHERE post_id = ? AND user_id = ?", post.ID, userID).Scan(&post.UserVote)
		if err == sql.ErrNoRows {
			post.UserVote = 0
		} else if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) GetCommentsByUserIDForPost(postID, userID int) ([]*Comment, error) {
	stmt := `
		SELECT comments.id, comments.post_id, comments.user_id, comments.created, comments.content
		FROM comments
		WHERE comments.post_id = ? AND comments.user_id = ?
		ORDER BY comments.created ASC
	`
	rows, err := m.DB.Query(stmt, postID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Created, &comment.Content)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
