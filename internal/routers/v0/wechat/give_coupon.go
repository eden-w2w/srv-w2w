package wechat

import (
	"context"
	"fmt"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-w2w/lib-modules/modules/id_generator"
	"github.com/eden-w2w/lib-modules/modules/wechat"
	"github.com/eden-w2w/srv-w2w/internal/constants/errors"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
	"github.com/eden-w2w/wechatpay-go/services/coupons"
)

func init() {
	Router.Register(courier.NewRouter(middleware.Authorization{}, GiveCoupon{}))
}

// GiveCoupon 发放代金券
type GiveCoupon struct {
	httpx.MethodGet
	StockID string `in:"path" name:"stockID"`
}

func (req GiveCoupon) Path() string {
	return "/stocks/:stockID/coupon"
}

func (req GiveCoupon) Output(ctx context.Context) (result interface{}, err error) {
	u := middleware.GetUserByContext(ctx)
	if u == nil {
		return nil, errors.Unauthorized
	}

	id, _ := id_generator.GetGenerator().GenerateUniqueID()
	_, err = wechat.GetController().GiveCoupon(coupons.GiveCouponRequest{
		StockID:           req.StockID,
		OpenID:            &u.OpenID,
		OutRequestNo:      fmt.Sprintf("%d", id),
		AppID:             global.Config.Wechat.AppID,
		StockCreatorMchID: global.Config.Wechat.MerchantID,
	})
	return
}
