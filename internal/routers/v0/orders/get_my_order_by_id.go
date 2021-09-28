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
	Router.Register(courier.NewRouter(GetMyOrderByID{}))
}

// GetMyOrderByID 获取一条订单
type GetMyOrderByID struct {
	httpx.MethodGet
	// 订单ID
	OrderID uint64 `in:"path" name:"orderID"`
}

func (req GetMyOrderByID) Path() string {
	return "/:orderID"
}

func (req GetMyOrderByID) Output(ctx context.Context) (result interface{}, err error) {
	user := middleware.GetUserByContext(ctx)
	if user == nil {
		return nil, errors.Unauthorized
	}
	o, err := order.GetController().GetOrder(req.OrderID, user.UserID, nil, false)
	if err != nil {
		return nil, err
	}

	response := &order.GetMyOrderByIDResponse{
		OrderID:       o.OrderID,
		UserID:        o.UserID,
		RefererID:     o.RefererID,
		TotalPrice:    o.TotalPrice,
		PaymentMethod: o.PaymentMethod,
		Remark:        o.Remark,
		Recipients:    o.Recipients,
		ShippingAddr:  o.ShippingAddr,
		Mobile:        o.Mobile,
		Status:        o.Status,
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
		Goods:         make([]order.OrderGoodsListResponse, 0),
	}
	goods, err := order.GetController().GetOrderGoods(o.OrderID)
	if err != nil {
		return nil, err
	}
	for _, g := range goods {
		response.Goods = append(response.Goods, order.OrderGoodsListResponse{
			GoodsID:        g.GoodsID,
			Name:           g.Name,
			MainPicture:    g.MainPicture,
			Specifications: g.Specifications,
			Price:          g.Price,
			Amount:         g.Amount,
		})
	}
	return response, nil
}
