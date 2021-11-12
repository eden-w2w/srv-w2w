package wechat

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-w2w/lib-modules/modules/wechat"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
	"github.com/eden-w2w/wechatpay-go/core"
	"github.com/eden-w2w/wechatpay-go/services/coupons"
)

func init() {
	Router.Register(courier.NewRouter(middleware.Authorization{}, GetCouponStocks{}))
}

// GetCouponStocks 获取可用的代金券批次列表
type GetCouponStocks struct {
	httpx.MethodGet
}

func (req GetCouponStocks) Path() string {
	return "/get_coupon_stocks"
}

type GetCouponStocksResponse struct {
	Data  []coupons.Stock `json:"data"`
	Total int64           `json:"total"`
}

func (req GetCouponStocks) Output(ctx context.Context) (result interface{}, err error) {
	request := coupons.GetStocksRequest{
		Offset:            0,
		Limit:             10,
		StockCreatorMchID: global.Config.Wechat.MerchantID,
		CreateStartTime:   nil,
		CreateEndTime:     nil,
		Status:            core.String("running"),
	}
	response, err := wechat.GetController().GetStocks(request)
	if err != nil {
		return
	}

	return &GetCouponStocksResponse{
		Data:  response.Data,
		Total: response.Total,
	}, nil
}
