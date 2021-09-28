package users

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	errors "github.com/eden-w2w/lib-modules/constants/general_errors"
	"github.com/eden-w2w/lib-modules/modules/user"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
)

func init() {
	Router.Register(courier.NewRouter(UpdateUserInfo{}))
}

// UpdateUserInfo 更新用户信息
type UpdateUserInfo struct {
	httpx.MethodPut
	Body user.UpdateUserInfoParams `in:"body"`
}

func (req UpdateUserInfo) Path() string {
	return "/updateUserInfo"
}

func (req UpdateUserInfo) Output(ctx context.Context) (result interface{}, err error) {
	u := middleware.GetUserByContext(ctx)
	if u == nil {
		return nil, errors.Unauthorized
	}

	err = user.GetController().UpdateUserInfo(u.UserID, req.Body)
	return
}
