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
	Router.Register(courier.NewRouter(CreateShippingAddress{}))
}

// CreateShippingAddress 创建收货地址
type CreateShippingAddress struct {
	httpx.MethodPost
	Data user.CreateShippingAddressParams `in:"body"`
}

func (req CreateShippingAddress) Path() string {
	return "/address"
}

func (req CreateShippingAddress) Output(ctx context.Context) (result interface{}, err error) {
	u := middleware.GetUserByContext(ctx)
	if u == nil {
		return nil, errors.Unauthorized
	}

	req.Data.UserID = u.UserID
	return user.GetController().CreateShippingAddress(req.Data, nil)
}
