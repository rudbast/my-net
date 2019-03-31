package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/rudbast/my-net/util"
)

func (m *Module) HandleUserPost(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()

	uid := util.GetContextRequester(ctx)

	var form struct {
		Content string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		return nil, errors.Wrap(err, "handler/post/param")
	}

	return nil, m.Service.AddPost(ctx, uid, form.Content, time.Now())
}

func (m *Module) HandleUserFeed(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()

	uid := util.GetContextRequester(ctx)
	// Query posts in the last X hour.
	// TODO: Make it configurable ?
	date := time.Now().Add(time.Hour * -36)

	return m.Service.QueryRecentPosts(ctx, uid, date)
}
