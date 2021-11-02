package wechat

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-framework/courier/transport_http"
	"github.com/eden-framework/sqlx"
	errors "github.com/eden-w2w/lib-modules/constants/general_errors"
	"github.com/eden-w2w/lib-modules/modules/wechat"
	"github.com/eden-w2w/srv-w2w/internal/global"
	wechatModule "github.com/eden-w2w/srv-w2w/internal/modules/wechat"
)

func init() {
	Router.Register(courier.NewRouter(PaymentNotify{}))
}

// PaymentNotify 微信支付回调
type PaymentNotify struct {
	httpx.MethodPost
}

func (req PaymentNotify) Path() string {
	return "/payment_notify"
}

func (req PaymentNotify) Output(ctx context.Context) (result interface{}, err error) {
	request := transport_http.GetRequest(ctx)
	_, trans, err := wechat.GetController().ParseWechatPaymentNotify(ctx, request)
	if err != nil {
		return nil, err
	}

	tx := sqlx.NewTasks(global.Config.MasterDB)
	tx = tx.With(
		func(db sqlx.DBExecutor) error {
			return wechatModule.UpdatePaymentByWechat(trans, db)
		})

	err = tx.Do()
	if err != nil {
		return nil, errors.InternalError
	}
	return wechat.WechatNotifyResponse{
		Code:    "SUCCESS",
		Message: "",
	}, nil
}
