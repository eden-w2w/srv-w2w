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
	o, l, err := order.GetController().GetOrder(req.OrderID, user.UserID, nil, false)
	if err != nil {
		return nil, err
	}

	response := &order.GetOrderByIDResponse{
		OrderID:        o.OrderID,
		UserID:         o.UserID,
		NickName:       o.NickName,
		OpenID:         o.UserOpenID,
		TotalPrice:     o.TotalPrice,
		DiscountAmount: o.DiscountAmount,
		ActualAmount:   o.ActualAmount,
		PaymentMethod:  o.PaymentMethod,
		Remark:         o.Remark,
		Recipients:     l.Recipients,
		ShippingAddr:   l.ShippingAddr,
		Mobile:         l.Mobile,
		CourierCompany: l.CourierCompany,
		CourierNumber:  l.CourierNumber,
		Status:         o.Status,
		ExpiredAt:      o.ExpiredAt,
		CreatedAt:      o.CreatedAt,
		UpdatedAt:      o.UpdatedAt,
		Goods:          make([]order.GoodsListResponse, 0),
	}
	goods, err := order.GetController().GetOrderGoods(o.OrderID)
	if err != nil {
		return nil, err
	}
	for _, g := range goods {
		response.Goods = append(response.Goods, order.GoodsListResponse{
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
