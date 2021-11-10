package payment

import (
	"context"
	"fmt"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-framework/sqlx"
	"github.com/eden-w2w/lib-modules/constants/enums"
	"github.com/eden-w2w/lib-modules/databases"
	"github.com/eden-w2w/lib-modules/modules/order"
	"github.com/eden-w2w/lib-modules/modules/payment_flow"
	"github.com/eden-w2w/lib-modules/modules/wechat"
	"github.com/eden-w2w/srv-w2w/internal/constants/errors"
	"github.com/eden-w2w/srv-w2w/internal/global"
	wechatModule "github.com/eden-w2w/srv-w2w/internal/modules/wechat"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
	"github.com/eden-w2w/wechatpay-go/core"
	"github.com/eden-w2w/wechatpay-go/services/payments/jsapi"
	"github.com/sirupsen/logrus"
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
	var retErr error
	tx = tx.With(
		func(db sqlx.DBExecutor) error {
			model, _, err := order.GetController().GetOrder(req.Data.OrderID, user.UserID, db, true)
			if err != nil {
				return err
			}
			o = model

			flows, err := payment_flow.GetController().GetFlowByOrderIDAndStatus(
				o.OrderID,
				user.UserID,
				[]enums.PaymentStatus{enums.PAYMENT_STATUS__CREATED, enums.PAYMENT_STATUS__PROCESS},
				db,
			)
			if err != nil {
				return err
			}
			for _, flow := range flows {
				err = wechat.GetController().CloseOrder(
					jsapi.CloseOrderRequest{
						OutTradeNo: core.String(fmt.Sprintf("%d", flow.FlowID)),
						Mchid:      core.String(global.Config.Wechat.MerchantID),
					},
				)
				if err != nil {
					// 查单并更新支付单状态
					tran, err := wechat.GetController().QueryOrderByOutTradeNo(
						jsapi.QueryOrderByOutTradeNoRequest{
							OutTradeNo: core.String(fmt.Sprintf("%d", flow.FlowID)),
							Mchid:      core.String(global.Config.Wechat.MerchantID),
						},
					)
					if err != nil {
						return errors.LastPaymentCloseConflict
					}
					_ = wechatModule.UpdatePaymentByWechat(tran, nil)
					return errors.LastPaymentCloseConflict
				}
			}
			return nil
		},
	)

	var paymentFlow *databases.PaymentFlow
	tx = tx.With(
		func(db sqlx.DBExecutor) error {
			req.Data.UserID = user.UserID
			req.Data.Amount = o.ActualAmount
			paymentFlow, err = payment_flow.GetController().CreatePaymentFlow(req.Data, db)
			return err
		},
	)

	var wechatResp *jsapi.PrepayWithRequestPaymentResponse
	tx = tx.With(
		func(_ sqlx.DBExecutor) error {
			resp, err := wechat.GetController().CreatePrePayment(ctx, o, paymentFlow, user)
			if err != nil {
				return err
			}
			wechatResp = resp
			return nil
		},
	)

	tx = tx.With(
		func(db sqlx.DBExecutor) error {
			if wechatResp == nil {
				return nil
			}
			return payment_flow.GetController().UpdatePaymentFlowRemoteID(paymentFlow.FlowID, *wechatResp.PrepayId, db)
		},
	)

	err = tx.Do()
	if err != nil {
		logrus.Errorf("[CreatePaymentFlow] tx.Do() err: %v, params: %+v", err, req.Data)
		return
	}
	if retErr != nil {
		return nil, retErr
	}
	return payment_flow.CreatePaymentFlowResponse{
		PaymentFlow:  paymentFlow,
		WechatPrepay: wechatResp,
	}, nil
}
