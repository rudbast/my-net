package core

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rudbast/my-net/model"
)

func (s *Service) AddFollow(ctx context.Context, uid, followedID int64) error {
	err := model.InsertFollow(ctx, s.database, model.Follow{
		UserID:     uid,
		FollowedID: followedID,
	})
	if err != nil {
		return errors.Wrap(err, "service/follow")
	}

	return nil
}

func (s *Service) RemoveFollow(ctx context.Context, uid, followedID int64) error {
	err := model.DeleteFollow(ctx, s.database, model.Follow{
		UserID:     uid,
		FollowedID: followedID,
	})
	if err != nil {
		return errors.Wrap(err, "service/unfollow")
	}

	return nil
}
