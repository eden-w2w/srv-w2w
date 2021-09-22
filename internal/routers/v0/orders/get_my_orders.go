package orders

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-w2w/srv-w2w/internal/contants/errors"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"

	"github.com/eden-w2w/srv-w2w/internal/modules/order"
)

func init() {
	Router.Register(courier.NewRouter(GetMyOrders{}))
}

// GetMyOrders 获取我的订单列表
type GetMyOrders struct {
	httpx.MethodGet

	order.GetOrdersParams
}

func (req GetMyOrders) Path() string {
	return ""
}

func (req GetMyOrders) Output(ctx context.Context) (result interface{}, err error) {
	user := middleware.GetUserByContext(ctx)
	if user == nil {
		return nil, errors.Unauthorized
	}

	req.UserID = user.UserID
	orders, err := order.GetController().GetOrders(req.GetOrdersParams)
	if err != nil {
		return nil, err
	}

	var data = make([]order.GetMyOrdersResponse, 0)
	for _, o := range orders {
		orderResp := order.GetMyOrdersResponse{
			OrderID:       o.OrderID,
			UserID:        o.UserID,
			TotalPrice:    o.TotalPrice,
			PaymentMethod: o.PaymentMethod,
			Status:        o.Status,
			CreatedAt:     o.CreatedAt,
			Goods:         make([]order.OrderGoodsListResponse, 0),
		}
		goods, err := order.GetController().GetOrderGoods(o.OrderID)
		if err != nil {
			return nil, err
		}
		for _, g := range goods {
			orderResp.Goods = append(orderResp.Goods, order.OrderGoodsListResponse{
				GoodsID:        g.GoodsID,
				Name:           g.Name,
				MainPicture:    g.MainPicture,
				Specifications: g.Specifications,
				Price:          g.Price,
				Amount:         g.Amount,
			})
		}
		data = append(data, orderResp)
	}
	return data, nil
}