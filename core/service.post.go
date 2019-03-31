package core

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rudbast/my-net/model"
)

func (s *Service) AddPost(ctx context.Context, uid int64, content string, postedAt time.Time) error {
	err := model.InsertPost(ctx, s.database, model.Post{
		UserID:   uid,
		Content:  content,
		PostedAt: postedAt,
	})
	if err != nil {
		return errors.Wrap(err, "service/post/add")
	}

	return nil
}

func (s *Service) QueryRecentPosts(ctx context.Context, uid int64, recent time.Time) ([]model.Post, error) {
	posts, err := model.QueryPostsByUserFollowedFilterDate(ctx, s.database, uid, recent)
	if err != nil {
		return nil, errors.Wrap(err, "service/post/recent")
	}

	return posts, nil
}
