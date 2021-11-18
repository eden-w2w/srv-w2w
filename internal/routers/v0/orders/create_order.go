package orders

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	errors "github.com/eden-w2w/lib-modules/constants/general_errors"
	"github.com/eden-w2w/lib-modules/modules/order"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
)

func init() {
	Router.Register(courier.NewRouter(CreateOrder{}))
}

// CreateOrder 创建订单
type CreateOrder struct {
	httpx.MethodPost
	Data order.CreateOrderParams `in:"body"`
}

func (req CreateOrder) Path() string {
	return ""
}

func (req CreateOrder) Output(ctx context.Context) (result interface{}, err error) {
	user := middleware.GetUserByContext(ctx)
	if user == nil {
		return nil, errors.Unauthorized
	}

	req.Data.UserID = user.UserID
	return order.GetController().CreateOrder(req.Data)
}
