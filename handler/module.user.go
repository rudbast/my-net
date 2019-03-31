package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rudbast/my-net/util"
)

func (m *Module) HandleUserFollow(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	vars := mux.Vars(r)

	uid := util.GetContextRequester(ctx)
	followedID, err := strconv.ParseInt(vars["followed_id"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "handler/follow/param")
	}

	return nil, m.Service.AddFollow(ctx, uid, followedID)
}

func (m *Module) HandleUserUnfollow(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	vars := mux.Vars(r)

	uid := util.GetContextRequester(ctx)
	followedID, err := strconv.ParseInt(vars["followed_id"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "handler/unfollow/param")
	}

	return nil, m.Service.RemoveFollow(ctx, uid, followedID)
}
