package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

type (
	Post struct {
		ID       int64     `json:"id"`
		UserID   int64     `json:"user_id"`
		Content  string    `json:"content"`
		PostedAt time.Time `json:"posted_at"`
	}
)

func InsertPost(ctx context.Context, db *sql.DB, post Post) error {
	query := "INSERT INTO posts (user_id, content, posted_at) VALUES ($1, $2, $3)"

	_, err := db.ExecContext(ctx, query, post.UserID, post.Content, post.PostedAt)
	if err != nil {
		return errors.Wrap(err, "model/post/insert")
	}

	return nil
}

// TODO: Add pagination.
func QueryPostsByUserFollowedFilterDate(ctx context.Context, db *sql.DB, userID int64, recent time.Time) ([]Post, error) {
	query := `
		WITH cte_followed AS (
			SELECT
				f.followed_id
			FROM
				users u
			INNER JOIN
				follows f ON u.id = f.user_id
			WHERE
				u.id = $1
		)
		SELECT
			p.id,
			p.user_id,
			p.content,
			p.posted_at
		FROM
			posts p
		INNER JOIN
			cte_followed c ON p.user_id = c.followed_id
		WHERE
			p.posted_at > $2
		ORDER BY
			p.posted_at DESC
	`

	rows, err := db.QueryContext(ctx, query, userID, recent)
	if err != nil {
		return nil, errors.Wrap(err, "model/post/followed")
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post

		err = rows.Scan(&post.ID, &post.UserID, &post.Content, &post.PostedAt)
		if err != nil {
			return nil, errors.Wrap(err, "model/post/followed")
		}

		posts = append(posts, post)
	}

	return posts, nil
}
