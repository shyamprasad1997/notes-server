package utils

import (
	"context"
	"notes-server/constants"

	"github.com/google/uuid"
)

func GetRequestIDFromCtx(ctx context.Context) string {
	if id, ok := ctx.Value(constants.RequestIDCtxKey).(string); ok {
		return id
	}
	return ""
}

func GetEmailFromCtx(ctx context.Context) string {
	if email, ok := ctx.Value(constants.EmailCtxKey).(string); ok {
		return email
	}
	return ""
}

func NewID() int32 {
	u, _ := uuid.NewRandom()
	return int32(u.ID())
}
