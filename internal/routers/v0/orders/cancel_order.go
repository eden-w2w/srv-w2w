package orders

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	errors "github.com/eden-w2w/lib-modules/constants/general_errors"
	"github.com/eden-w2w/lib-modules/modules/goods"
	"github.com/eden-w2w/lib-modules/modules/order"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
)

func init() {
	Router.Register(courier.NewRouter(CancelOrder{}))
}

// CancelOrder 取消订单
type CancelOrder struct {
	httpx.MethodDelete
	// 订单ID
	OrderID uint64 `in:"path" name:"orderID,string"`
}

func (req CancelOrder) Path() string {
	return "/:orderID"
}

func (req CancelOrder) Output(ctx context.Context) (result interface{}, err error) {
	user := middleware.GetUserByContext(ctx)
	if user == nil {
		return nil, errors.Unauthorized
	}

	err = order.GetController().CancelOrder(req.OrderID, user.UserID, goods.GetController().UnlockInventory)
	return
}
