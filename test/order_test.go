package test

import (
	"context"
	"github.com/eden-w2w/lib-modules/constants/enums"
	errors "github.com/eden-w2w/lib-modules/constants/general_errors"
	"github.com/eden-w2w/lib-modules/databases"
	"github.com/eden-w2w/lib-modules/modules"
	"github.com/eden-w2w/lib-modules/modules/order"
	"github.com/eden-w2w/lib-modules/modules/promotion_flow"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
	"github.com/eden-w2w/srv-w2w/internal/routers/v0/orders"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testCreateOrder(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.AuthContextKey, orderUserModel)

	request := orders.CreateOrder{
		Data: order.CreateOrderParams{
			TotalPrice:    12000,
			PaymentMethod: enums.PAYMENT_METHOD__WECHAT,
			Remark:        "这是单元测试订单",
			Recipients:    "测试人员",
			ShippingAddr:  "测试工厂",
			Mobile:        "137********",
			Goods: []order.CreateOrderGoodsParams{
				{
					GoodsID: 10001,
					Amount:  1,
				},
			},
		},
	}
	orderResp, err := request.Output(ctx)
	assert.Nil(t, err)
	orderModel = orderResp.(*databases.Order)
	assert.Equal(t, uint64(12000), orderModel.TotalPrice)

	goods, err := order.GetController().GetOrderGoods(orderModel.OrderID)
	assert.Nil(t, err)
	assert.Len(t, goods, 1)
	assert.Equal(t, uint64(10001), goods[0].GoodsID)
	assert.Equal(t, uint32(1), goods[0].Amount)
}

func testUpdateOrderStatusIncorrect(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.AuthContextKey, orderUserModel)

	request := orders.TestUpdateOrderStatus{
		OrderID: orderModel.OrderID,
		Status:  enums.ORDER_STATUS__COMPLETE,
	}
	_, err := request.Output(ctx)
	assert.Equal(t, errors.OrderStatusFlowIncorrect, err)
}

func testUpdateOrderStatus(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.AuthContextKey, orderUserModel)

	request := orders.TestUpdateOrderStatus{
		OrderID: orderModel.OrderID,
		Status:  enums.ORDER_STATUS__CONFIRM,
	}
	_, err := request.Output(ctx)
	assert.Nil(t, err)

	request = orders.TestUpdateOrderStatus{
		OrderID: orderModel.OrderID,
		Status:  enums.ORDER_STATUS__DISPATCH,
	}
	_, err = request.Output(ctx)
	assert.Nil(t, err)

	request = orders.TestUpdateOrderStatus{
		OrderID: orderModel.OrderID,
		Status:  enums.ORDER_STATUS__COMPLETE,
	}
	_, err = request.Output(ctx)
	assert.Nil(t, err)

	promotionFlowModel, err = promotion_flow.GetController().GetPromotionFlows(promotion_flow.GetPromotionFlowParams{
		UserID: orderUserModel.RefererID,
		Pagination: modules.Pagination{
			Size: -1,
		},
	})
	assert.Nil(t, err)
	assert.Len(t, promotionFlowModel, 1)
}
