package db

import (
	"database/sql"
	"fmt"

	"github.com/yves1da2vinci3/go_elastic/models"
)

var (
	ErrNoRecord = fmt.Errorf("no matching record found")
	insertOp    = "insert"
	deleteOp    = "delete"
	updateOp    = "update"
)

func (db Database) SavePost(post *models.Post) error {
	var id int
	query := `INSERT INTO posts(title, body) VALUES ($1, $2) RETURNING id`
	err := db.Conn.QueryRow(query, post.Title, post.Body).Scan(&id)
	if err != nil {
		return err
	}

	// doing this at app level, but if you feel like database operations wont always pass through
	// the application, you can move it to the DB level using triggers.
	logQuery := `INSERT INTO post_logs(post_id, operation) VALUES ($1, $2)`
	post.ID = id
	_, err = db.Conn.Exec(logQuery, post.ID, insertOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) UpdatePost(postId int, post models.Post) error {
	query := "UPDATE posts SET title=$1, body=$2 WHERE id=$3"
	_, err := db.Conn.Exec(query, post.Title, post.Body, postId)
	if err != nil {
		return err
	}

	post.ID = postId
	logQuery := "INSERT INTO post_logs(post_id, operation) VALUES ($1, $2)"
	_, err = db.Conn.Exec(logQuery, post.ID, updateOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) DeletePost(postId int) error {
	query := "DELETE FROM Posts WHERE id=$1"
	_, err := db.Conn.Exec(query, postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRecord
		}
		return err
	}

	logQuery := "INSERT INTO post_logs(post_id, operation) VALUES ($1, $2)"
	_, err = db.Conn.Exec(logQuery, postId, deleteOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) GetPostById(postId int) (models.Post, error) {
	post := models.Post{}
	query := "SELECT id, title, body FROM posts WHERE id = $1"
	row := db.Conn.QueryRow(query, postId)
	switch err := row.Scan(&post.ID, &post.Title, &post.Body); err {
	case sql.ErrNoRows:
		return post, ErrNoRecord
	default:
		return post, err
	}
}

func (db Database) GetPosts() ([]models.Post, error) {
	var list []models.Post
	query := "SELECT id, title, body FROM posts ORDER BY id DESC"
	rows, err := db.Conn.Query(query)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Body)
		if err != nil {
			return list, err
		}
		list = append(list, post)
	}
	return list, nil
}
