package payment

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-framework/sqlx"
	errors "github.com/eden-w2w/lib-modules/constants/general_errors"
	"github.com/eden-w2w/lib-modules/databases"
	"github.com/eden-w2w/lib-modules/modules/order"
	"github.com/eden-w2w/lib-modules/modules/payment_flow"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/eden-w2w/srv-w2w/internal/modules/wechat"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
	"github.com/sirupsen/logrus"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
)

func init() {
	Router.Register(courier.NewRouter(CreatePaymentFlow{}))
}

// CreatePaymentFlow 创建支付流水单
type CreatePaymentFlow struct {
	httpx.MethodPost

	Data payment_flow.CreatePaymentFlowParams `in:"body"`
}

func (req CreatePaymentFlow) Path() string {
	return ""
}

func (req CreatePaymentFlow) Output(ctx context.Context) (result interface{}, err error) {
	user := middleware.GetUserByContext(ctx)
	if user == nil {
		return nil, errors.Unauthorized
	}

	tx := sqlx.NewTasks(global.Config.MasterDB)

	var o *databases.Order
	tx = tx.With(func(db sqlx.DBExecutor) error {
		model, err := order.GetController().GetOrder(req.Data.OrderID, user.UserID, db, true)
		if err != nil {
			return err
		}
		o = model
		return nil
	})

	var paymentFlow *databases.PaymentFlow
	tx = tx.With(func(db sqlx.DBExecutor) error {
		req.Data.UserID = user.UserID
		req.Data.Amount = o.TotalPrice
		paymentFlow, err = payment_flow.GetController().CreatePaymentFlow(req.Data, db)
		return err
	})

	var wechatResp *jsapi.PrepayWithRequestPaymentResponse
	tx = tx.With(func(_ sqlx.DBExecutor) error {
		resp, err := wechat.GetController().CreatePrePayment(ctx, o, paymentFlow, user)
		if err != nil {
			return err
		}
		wechatResp = resp
		return nil
	})

	tx = tx.With(func(db sqlx.DBExecutor) error {
		if wechatResp == nil {
			return nil
		}
		return payment_flow.GetController().UpdatePaymentFlowRemoteID(paymentFlow.FlowID, *wechatResp.PrepayId, db)
	})

	err = tx.Do()
	if err != nil {
		logrus.Errorf("[CreatePaymentFlow] tx.Do() err: %v, params: %+v", err, req.Data)
		return
	}
	return payment_flow.CreatePaymentFlowResponse{
		PaymentFlow:  paymentFlow,
		WechatPrepay: wechatResp,
	}, nil
}
