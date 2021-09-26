package test

import (
	"github.com/eden-w2w/srv-w2w/internal/databases"
	"github.com/eden-w2w/srv-w2w/internal/modules/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

var userModel *databases.User
var orderModel *databases.Order
var paymentFlowModel *databases.PaymentFlow

func TestAll(t *testing.T) {
	u, err := user.GetController().GetUserByUserID(1441156099381141504, nil, false)
	assert.Nil(t, err)
	userModel = u

	// 创建订单
	t.Run("testCreateOrder", testCreateOrder)

	// 支付
	t.Run("testCreatePaymentFlow", testCreatePaymentFlow)
	t.Run("testPaymentNotifySuccess", testPaymentNotifySuccess)

	// 错误的订单状态流转
	t.Run("testUpdateOrderStatusIncorrect", testUpdateOrderStatusIncorrect)
	// 正确的订单状态流转
	t.Run("testUpdateOrderStatus", testUpdateOrderStatus)
}
