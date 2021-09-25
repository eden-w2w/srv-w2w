package wechat

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-framework/courier/transport_http"
	"github.com/eden-framework/sqlx"
	"github.com/eden-w2w/srv-w2w/internal/contants/enums"
	"github.com/eden-w2w/srv-w2w/internal/contants/errors"
	"github.com/eden-w2w/srv-w2w/internal/databases"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/eden-w2w/srv-w2w/internal/modules/order"
	"github.com/eden-w2w/srv-w2w/internal/modules/payment_flow"
	"github.com/eden-w2w/srv-w2w/internal/modules/wechat"
	"github.com/sirupsen/logrus"
	"strconv"
)

func init() {
	Router.Register(courier.NewRouter(Notify{}))
}

// Notify 微信支付回调
type Notify struct {
	httpx.MethodPost
}

func (req Notify) Path() string {
	return "/notify"
}

func (req Notify) Output(ctx context.Context) (result interface{}, err error) {
	request := transport_http.GetRequest(ctx)
	_, trans, err := wechat.GetController().ParseWechatPaymentNotify(ctx, request)
	if err != nil {
		return nil, err
	}
	tradeState, err := enums.ParseWechatTradeStateFromString(*trans.TradeState)
	if err != nil {
		logrus.Errorf("[Notify] enums.ParseWechatTradeStateFromString err: %v, TradeState: %s", err, *trans.TradeState)
		return nil, errors.InternalError
	}
	if !tradeState.IsEnding() {
		logrus.Infof("[Notify] !tradeState.IsEnding(), state: %s", tradeState.String())
		return nil, nil
	}
	flowID, err := strconv.ParseUint(*trans.OutTradeNo, 10, 64)
	if err != nil {
		logrus.Errorf("[Notify] strconv.ParseUint err: %v, OutTradeNo: %s", err, *trans.OutTradeNo)
		return nil, errors.InternalError
	}
	amount := uint64(*trans.Amount.Total)

	tx := sqlx.NewTasks(global.Config.MasterDB)
	var paymentFlow *databases.PaymentFlow
	tx = tx.With(func(db sqlx.DBExecutor) error {
		paymentFlow, err = payment_flow.GetController().GetPaymentFlowByID(flowID, db, true)
		if err != nil {
			return err
		}

		if paymentFlow.Amount != amount {
			return errors.FlowAmountIncorrect
		}
		return nil
	})

	tx = tx.With(func(db sqlx.DBExecutor) (err error) {
		if tradeState.IsEqual(paymentFlow.Status) {
			return nil
		}
		if tradeState.IsSuccess() {
			err = payment_flow.GetController().UpdatePaymentFlowStatus(paymentFlow.FlowID, enums.PAYMENT_STATUS__SUCCESS, trans, db)
			if err != nil {
				return
			}
			// 联动订单
			err = order.GetController().UpdateOrderStatusWithDB(db, paymentFlow.OrderID, enums.ORDER_STATUS__PAID)
		} else if tradeState.IsFail() {
			err = payment_flow.GetController().UpdatePaymentFlowStatus(paymentFlow.FlowID, enums.PAYMENT_STATUS__FAIL, trans, db)
			if err != nil {
				return
			}
		}
		return
	})

	err = tx.Do()
	if err != nil {
		logrus.Errorf("[Notify] tx.Do() err: %v, trans: %+v", err, trans)
		return nil, err
	}
	return wechat.WechatNotifyResponse{
		Code:    "SUCCESS",
		Message: "",
	}, nil
}
