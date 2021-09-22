package middleware

import (
	"context"
	"github.com/eden-w2w/srv-w2w/internal/databases"
)

const AuthContextKey = "Authorization"

func GetUserByContext(ctx context.Context) *databases.User {
	val := ctx.Value(AuthContextKey)
	if user, ok := val.(*databases.User); ok {
		return user
	}
	return nil
}
