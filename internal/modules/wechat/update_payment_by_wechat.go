package wechat

import (
	"github.com/eden-framework/sqlx"
	"github.com/eden-w2w/lib-modules/constants/enums"
	errors "github.com/eden-w2w/lib-modules/constants/general_errors"
	"github.com/eden-w2w/lib-modules/databases"
	"github.com/eden-w2w/lib-modules/modules/goods"
	"github.com/eden-w2w/lib-modules/modules/order"
	"github.com/eden-w2w/lib-modules/modules/payment_flow"
	"github.com/sirupsen/logrus"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"strconv"
)

func UpdatePaymentByWechat(tran *payments.Transaction, db sqlx.DBExecutor) error {
	tradeState, err := enums.ParseWechatTradeStateFromString(*tran.TradeState)
	if err != nil {
		logrus.Errorf("[PaymentNotify] enums.ParseWechatTradeStateFromString err: %v, TradeState: %s", err, *tran.TradeState)
		return err
	}
	if !tradeState.IsEnding() {
		logrus.Infof("[PaymentNotify] !tradeState.IsEnding(), state: %s", tradeState.String())
		return nil
	}

	flowID, err := strconv.ParseUint(*tran.OutTradeNo, 10, 64)
	if err != nil {
		logrus.Errorf("[PaymentNotify] strconv.ParseUint err: %v, OutTradeNo: %s", err, *tran.OutTradeNo)
		return errors.InternalError
	}

	var paymentFlow *databases.PaymentFlow
	paymentFlow, err = payment_flow.GetController().GetPaymentFlowByID(flowID, db, true)
	if err != nil {
		return err
	}

	if !tradeState.IsFail() {
		amount := uint64(*tran.Amount.Total)
		if paymentFlow.Amount != amount {
			return errors.FlowAmountIncorrect
		}
	}
	if tradeState.IsEqual(paymentFlow.Status) {
		return nil
	}
	if tradeState.IsSuccess() {
		err = payment_flow.GetController().UpdatePaymentFlowStatus(
			paymentFlow,
			enums.PAYMENT_STATUS__SUCCESS,
			tran,
			db,
		)
		if err != nil {
			return err
		}
		// 联动订单
		var orderModel *databases.Order
		var logistics *databases.OrderLogistics
		orderModel, logistics, err = order.GetController().GetOrder(
			paymentFlow.OrderID,
			paymentFlow.UserID,
			db,
			true,
		)
		if err != nil {
			return err
		}
		orderGoods, err := order.GetController().GetOrderGoods(paymentFlow.OrderID, db)
		if err != nil {
			return err
		}
		err = order.GetController().UpdateOrder(
			orderModel, logistics, orderGoods, order.UpdateOrderParams{
				Status: enums.ORDER_STATUS__PAID,
			}, goods.GetController().LockInventory, goods.GetController().UnlockInventory, db,
		)
		if err != nil {
			return err
		}
	} else if tradeState.IsFail() {
		err = payment_flow.GetController().UpdatePaymentFlowStatus(
			paymentFlow,
			tradeState.ToPaymentStatus(),
			tran,
			db,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
