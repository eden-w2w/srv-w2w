package wechat

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-framework/sqlx/datatypes"
	"github.com/eden-w2w/lib-modules/modules/wechat"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
	"github.com/eden-w2w/wechatpay-go/core"
	"github.com/eden-w2w/wechatpay-go/services/coupons"
)

func init() {
	Router.Register(courier.NewRouter(middleware.Authorization{}, CheckCouponExist{}))
}

// CheckCouponExist 检查是否存在可用的代金券批次
type CheckCouponExist struct {
	httpx.MethodGet
}

func (req CheckCouponExist) Path() string {
	return "/check_coupons_exist"
}

type CheckCouponExistResponse struct {
	Exist datatypes.Bool `json:"exist"`
}

func (req CheckCouponExist) Output(ctx context.Context) (result interface{}, err error) {
	request := coupons.GetStocksRequest{
		Offset:            0,
		Limit:             1,
		StockCreatorMchID: global.Config.Wechat.MerchantID,
		CreateStartTime:   nil,
		CreateEndTime:     nil,
		Status:            core.String("running"),
	}
	response, err := wechat.GetController().GetStocks(request)
	if err != nil {
		return
	}

	if len(response.Data) > 0 {
		return &CheckCouponExistResponse{Exist: datatypes.BOOL_TRUE}, nil
	}

	return &CheckCouponExistResponse{Exist: datatypes.BOOL_FALSE}, nil
}
