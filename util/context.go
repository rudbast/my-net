package util

import "context"

const (
	ctxKeyRequester = "requester"
)

func SetContextRequester(ctx context.Context, uid int64) context.Context {
	return context.WithValue(ctx, ctxKeyRequester, uid)
}

func GetContextRequester(ctx context.Context) int64 {
	uid, _ := ctx.Value(ctxKeyRequester).(int64)
	return uid
}
