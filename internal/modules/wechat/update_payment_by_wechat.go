package wechat

import (
	"github.com/eden-framework/sqlx"
	"github.com/eden-w2w/lib-modules/constants/enums"
	errors "github.com/eden-w2w/lib-modules/constants/general_errors"
	"github.com/eden-w2w/lib-modules/databases"
	"github.com/eden-w2w/lib-modules/modules/goods"
	"github.com/eden-w2w/lib-modules/modules/order"
	"github.com/eden-w2w/lib-modules/modules/payment_flow"
	"github.com/eden-w2w/lib-modules/pkg/webhook"
	"github.com/eden-w2w/wechatpay-go/services/payments"
	"github.com/sirupsen/logrus"
	"strconv"
)

func UpdatePaymentByWechat(tran *payments.Transaction, db sqlx.DBExecutor) error {
	tradeState, err := enums.ParseWechatTradeStateFromString(*tran.TradeState)
	if err != nil {
		logrus.Errorf(
			"[PaymentNotify] enums.ParseWechatTradeStateFromString err: %v, TradeState: %s",
			err,
			*tran.TradeState,
		)
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
		// 检查代金券信息
		var discountAmount, actualAmount = uint64(0), uint64(0)
		if len(tran.PromotionDetail) > 0 {
			for _, detail := range tran.PromotionDetail {
				discountAmount += uint64(*detail.Amount)
			}
			actualAmount = paymentFlow.Amount - discountAmount
			err = payment_flow.GetController().UpdatePaymentFlowAmount(
				paymentFlow.FlowID,
				discountAmount,
				actualAmount,
				db,
			)
			if err != nil {
				return err
			}
		}

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

		updateParams := order.UpdateOrderParams{
			Status: enums.ORDER_STATUS__PAID,
		}
		if len(tran.PromotionDetail) > 0 {
			updateParams.DiscountAmount = discountAmount
		}
		err = order.GetController().UpdateOrder(
			orderModel,
			logistics,
			orderGoods,
			updateParams,
			goods.GetController().LockInventory,
			goods.GetController().UnlockInventory,
			db,
		)
		if err != nil {
			return err
		}
		go webhook.GetInstance().SendPayment(orderModel, paymentFlow, logistics)
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
	return payment_flow.GetController().UpdatePaymentFlowRemoteID(paymentFlow.FlowID, *tran.TransactionId, db)
}
