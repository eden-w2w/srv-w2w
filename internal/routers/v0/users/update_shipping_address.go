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
	Router.Register(courier.NewRouter(UpdateShippingAddress{}))
}

// UpdateShippingAddress 更新收货地址
type UpdateShippingAddress struct {
	httpx.MethodPut
	// 业务ID
	ShippingID uint64                           `in:"path" name:"shippingID,string"`
	Data       user.UpdateShippingAddressParams `in:"body"`
}

func (req UpdateShippingAddress) Path() string {
	return "/address/:shippingID"
}

func (req UpdateShippingAddress) Output(ctx context.Context) (result interface{}, err error) {
	u := middleware.GetUserByContext(ctx)
	if u == nil {
		return nil, errors.Unauthorized
	}

	req.Data.ShippingID = req.ShippingID
	err = user.GetController().UpdateShippingAddress(req.Data, u.UserID, nil)
	return
}
