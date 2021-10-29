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
	Router.Register(courier.NewRouter(SetDefaultAddress{}))
}

type SetDefaultAddressRequest struct {
	ShippingID uint64 `json:"shippingID,string"`
}

// SetDefaultAddress 设置默认收货地址
type SetDefaultAddress struct {
	httpx.MethodPost
	Data SetDefaultAddressRequest `in:"body"`
}

func (req SetDefaultAddress) Path() string {
	return "/set_default_address"
}

func (req SetDefaultAddress) Output(ctx context.Context) (result interface{}, err error) {
	u := middleware.GetUserByContext(ctx)
	if u == nil {
		return nil, errors.Unauthorized
	}

	err = user.GetController().SetDefaultShippingAddress(u.UserID, req.Data.ShippingID, nil)
	return
}
