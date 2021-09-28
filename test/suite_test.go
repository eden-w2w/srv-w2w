package test

import (
	"github.com/eden-w2w/lib-modules/databases"
	"github.com/eden-w2w/lib-modules/modules/order"
	"github.com/eden-w2w/lib-modules/modules/user"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/stretchr/testify/assert"
	"testing"
)

var orderUserModel *databases.User
var promotionUserModel *databases.User
var orderModel *databases.Order
var logisticsModel *databases.OrderLogistics
var paymentFlowModel *databases.PaymentFlow
var promotionFlowModel []databases.PromotionFlow

func TestAll(t *testing.T) {
	u, err := user.GetController().GetUserByUserID(1441156099381141504, nil, false)
	assert.Nil(t, err)
	orderUserModel = u

	u, err = user.GetController().GetUserByUserID(1442548676592422912, nil, false)
	assert.Nil(t, err)
	promotionUserModel = u

	// 创建订单
	t.Run("testCreateOrder", testCreateOrder)

	// 支付
	t.Run("testCreatePaymentFlow", testCreatePaymentFlow)
	t.Run("testPaymentNotifySuccess", testPaymentNotifySuccess)

	// 错误的订单状态流转
	t.Run("testUpdateOrderStatusIncorrect", testUpdateOrderStatusIncorrect)
	// 正确的订单状态流转
	t.Run("testUpdateOrderStatus", testUpdateOrderStatus)

	// 查看我的钱包概览
	t.Run("testGetMyPromotionSummary", testGetMyPromotionSummary)

	// 清理
	goods, err := order.GetController().GetOrderGoods(orderModel.OrderID)
	assert.Nil(t, err)
	for _, good := range goods {
		_ = good.DeleteByOrderIDAndGoodsID(global.Config.MasterDB)
	}
	if orderModel != nil {
		_ = orderModel.DeleteByOrderID(global.Config.MasterDB)
	}
	if logisticsModel != nil {
		_ = logisticsModel.DeleteByLogisticsID(global.Config.MasterDB)
	}
	if paymentFlowModel != nil {
		_ = paymentFlowModel.DeleteByFlowID(global.Config.MasterDB)
	}
	for _, promotionFlow := range promotionFlowModel {
		_ = promotionFlow.DeleteByFlowID(global.Config.MasterDB)
	}
}
