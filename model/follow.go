package model

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type (
	Follow struct {
		ID         int64
		UserID     int64
		FollowedID int64
	}
)

func InsertFollow(ctx context.Context, db *sql.DB, follow Follow) error {
	query := "INSERT INTO follows (user_id, followed_id) VALUES ($1, $2)"

	_, err := db.ExecContext(ctx, query, follow.UserID, follow.FollowedID)
	if err != nil {
		return errors.Wrap(err, "model/follow/insert")
	}

	return nil
}

func DeleteFollow(ctx context.Context, db *sql.DB, follow Follow) error {
	query := "DELETE FROM follows WHERE user_id = $1 AND followed_id = $2"

	_, err := db.ExecContext(ctx, query, follow.UserID, follow.FollowedID)
	if err != nil {
		return errors.Wrap(err, "model/follow/delete")
	}

	return nil
}
