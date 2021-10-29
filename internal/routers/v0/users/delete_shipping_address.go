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
	Router.Register(courier.NewRouter(DeleteShippingAddress{}))
}

// DeleteShippingAddress 删除收货地址
type DeleteShippingAddress struct {
	httpx.MethodDelete
	// 业务ID
	ShippingID uint64 `in:"path" name:"shippingID,string"`
}

func (req DeleteShippingAddress) Path() string {
	return "/address/:shippingID"
}

func (req DeleteShippingAddress) Output(ctx context.Context) (result interface{}, err error) {
	u := middleware.GetUserByContext(ctx)
	if u == nil {
		return nil, errors.Unauthorized
	}

	err = user.GetController().DeleteShippingAddress(req.ShippingID, u.UserID, nil)
	return
}
