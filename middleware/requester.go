package middleware

import (
	"net/http"
	"strconv"

	"github.com/rudbast/my-net/util"
)

func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		uid, _ := strconv.ParseInt(r.Header.Get("UserID"), 10, 64)
		ctx = util.SetContextRequester(ctx, uid)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
