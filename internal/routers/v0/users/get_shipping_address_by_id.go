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
	Router.Register(courier.NewRouter(GetShippingAddressByID{}))
}

// GetShippingAddressByID 根据ID获取收货地址信息
type GetShippingAddressByID struct {
	httpx.MethodGet
	// 业务ID
	ShippingID uint64 `in:"path" name:"shippingID,string"`
}

func (req GetShippingAddressByID) Path() string {
	return "/address/:shippingID"
}

func (req GetShippingAddressByID) Output(ctx context.Context) (result interface{}, err error) {
	u := middleware.GetUserByContext(ctx)
	if u == nil {
		return nil, errors.Unauthorized
	}

	return user.GetController().GetShippingAddressByShippingID(req.ShippingID, u.UserID)
}
