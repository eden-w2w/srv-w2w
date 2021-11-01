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
	orders, _, err := order.GetController().GetOrders(req.GetOrdersParams, false)
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
			Goods:         make([]order.GoodsListResponse, 0),
		}
		goods, err := order.GetController().GetOrderGoods(o.OrderID, nil)
		if err != nil {
			return nil, err
		}
		for _, g := range goods {
			orderResp.Goods = append(orderResp.Goods, order.GoodsListResponse{
				GoodsID:        g.GoodsID,
				Name:           g.Name,
				Comment:        g.Comment,
				MainPicture:    g.MainPicture,
				Specifications: g.Specifications,
				Activities:     g.Activities,
				Price:          g.Price,
				Amount:         g.Amount,
			})
		}
		data = append(data, orderResp)
	}
	return data, nil
}
