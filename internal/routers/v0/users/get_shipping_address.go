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
	Router.Register(courier.NewRouter(GetShippingAddress{}))
}

// GetShippingAddress 获取收货地址
type GetShippingAddress struct {
	httpx.MethodGet
}

func (req GetShippingAddress) Path() string {
	return "/address"
}

func (req GetShippingAddress) Output(ctx context.Context) (result interface{}, err error) {
	u := middleware.GetUserByContext(ctx)
	if u == nil {
		return nil, errors.Unauthorized
	}

	return user.GetController().GetShippingAddressByUserID(u.UserID)
}
