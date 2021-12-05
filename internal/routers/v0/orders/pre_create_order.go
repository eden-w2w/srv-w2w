package orders

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-w2w/lib-modules/modules/order"
	"github.com/eden-w2w/srv-w2w/internal/constants/errors"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
)

func init() {
	Router.Register(courier.NewRouter(PreCreateOrder{}))
}

// PreCreateOrder 预创建订单，获取金额信息
type PreCreateOrder struct {
	httpx.MethodPost
	Data order.PreCreateOrderParams `in:"body"`
}

func (req PreCreateOrder) Path() string {
	return "/:orderID/pre"
}

type PreCreateOrderResponse struct {
	PreGoodsList  []order.PreCreateOrderGoodsParams `json:"preGoodsList"`
	TotalPrice    uint64                            `json:"totalPrice"`
	FreightName   string                            `json:"freightName"`
	FreightPrice  uint64                            `json:"freightPrice"`
	DiscountPrice uint64                            `json:"discountPrice"`
	ActualPrice   uint64                            `json:"actualPrice"`
}

func (req PreCreateOrder) Output(ctx context.Context) (result interface{}, err error) {
	user := middleware.GetUserByContext(ctx)
	if user == nil {
		return nil, errors.Unauthorized
	}

	req.Data.UserID = user.UserID
	preGoodsList, totalPrice, freightPrice, discountPrice, actualPrice, freightName, err := order.GetController().PreCreateOrder(req.Data)
	if err != nil {
		return
	}

	return &PreCreateOrderResponse{
		PreGoodsList:  preGoodsList,
		TotalPrice:    totalPrice,
		FreightName:   freightName,
		FreightPrice:  freightPrice,
		DiscountPrice: discountPrice,
		ActualPrice:   actualPrice,
	}, nil
}
