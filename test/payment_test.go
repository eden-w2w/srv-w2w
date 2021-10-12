package test

import (
	"context"
	"github.com/eden-w2w/lib-modules/constants/enums"
	"github.com/eden-w2w/lib-modules/modules/payment_flow"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
	"github.com/eden-w2w/srv-w2w/internal/routers/v0/payment"
	"github.com/eden-w2w/srv-w2w/internal/routers/v0/wechat"
	"github.com/stretchr/testify/require"
	"testing"
)

func testCreatePaymentFlow(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.AuthContextKey, orderUserModel)

	request := payment.CreatePaymentFlow{
		Data: payment_flow.CreatePaymentFlowParams{
			OrderID:       orderModel.OrderID,
			PaymentMethod: enums.PAYMENT_METHOD__WECHAT,
		},
	}
	resp, err := request.Output(ctx)
	require.Nil(t, err)
	response := resp.(payment_flow.CreatePaymentFlowResponse)
	paymentFlowModel = response.PaymentFlow

	require.Equal(t, orderModel.ActualAmount, paymentFlowModel.Amount)
}

func testPaymentNotifySuccess(t *testing.T) {
	ctx := context.Background()
	request := wechat.TestPaymentComplete{
		Amount:     paymentFlowModel.Amount,
		FlowID:     paymentFlowModel.FlowID,
		TradeState: enums.WECHAT_TRADE_STATE__SUCCESS,
	}
	_, err := request.Output(ctx)
	require.Nil(t, err)

	err = paymentFlowModel.FetchByFlowID(global.Config.MasterDB)
	require.Nil(t, err)
	require.Equal(t, enums.PAYMENT_STATUS__SUCCESS, paymentFlowModel.Status)

	err = orderModel.FetchByOrderID(global.Config.MasterDB)
	require.Nil(t, err)
	require.Equal(t, enums.ORDER_STATUS__PAID, orderModel.Status)
}
