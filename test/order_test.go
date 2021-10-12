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
	"github.com/stretchr/testify/require"
	"testing"
)

func testCreateOrder(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.AuthContextKey, orderUserModel)

	request := orders.CreateOrder{
		Data: order.CreateOrderParams{
			TotalPrice:     12000,
			DiscountAmount: 0,
			ActualAmount:   12000,
			PaymentMethod:  enums.PAYMENT_METHOD__WECHAT,
			Remark:         "这是单元测试订单",
			Recipients:     "测试人员",
			ShippingAddr:   "测试工厂",
			Mobile:         "137********",
			Goods: []order.CreateOrderGoodsParams{
				{
					GoodsID: 10001,
					Amount:  1,
				},
			},
		},
	}
	orderResp, err := request.Output(ctx)
	require.Nil(t, err)
	orderModel = orderResp.(*databases.Order)
	require.Equal(t, uint64(12000), orderModel.TotalPrice)
	require.Equal(t, uint64(0), orderModel.DiscountAmount)
	require.Equal(t, uint64(12000), orderModel.ActualAmount)

	_, logisticsModel, err = order.GetController().GetOrder(orderModel.OrderID, orderModel.UserID, nil, false)
	require.Nil(t, err)
	require.Equal(t, "测试人员", logisticsModel.Recipients)
	require.Equal(t, "测试工厂", logisticsModel.ShippingAddr)
	require.Equal(t, "137********", logisticsModel.Mobile)

	goods, err := order.GetController().GetOrderGoods(orderModel.OrderID)
	require.Nil(t, err)
	require.Len(t, goods, 1)
	require.Equal(t, uint64(10001), goods[0].GoodsID)
	require.Equal(t, uint32(1), goods[0].Amount)
}

func testUpdateOrderStatusIncorrect(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.AuthContextKey, orderUserModel)

	request := orders.TestUpdateOrderStatus{
		OrderID: orderModel.OrderID,
		Status:  enums.ORDER_STATUS__COMPLETE,
	}
	_, err := request.Output(ctx)
	require.Equal(t, errors.OrderStatusFlowIncorrect, err)
}

func testUpdateOrderStatus(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.AuthContextKey, orderUserModel)

	request := orders.TestUpdateOrderStatus{
		OrderID: orderModel.OrderID,
		Status:  enums.ORDER_STATUS__CONFIRM,
	}
	_, err := request.Output(ctx)
	require.Nil(t, err)

	request = orders.TestUpdateOrderStatus{
		OrderID: orderModel.OrderID,
		Status:  enums.ORDER_STATUS__DISPATCH,
	}
	_, err = request.Output(ctx)
	require.Nil(t, err)

	request = orders.TestUpdateOrderStatus{
		OrderID: orderModel.OrderID,
		Status:  enums.ORDER_STATUS__COMPLETE,
	}
	_, err = request.Output(ctx)
	require.Nil(t, err)

	promotionFlowModel, err = promotion_flow.GetController().GetPromotionFlows(promotion_flow.GetPromotionFlowParams{
		UserID: orderUserModel.RefererID,
		Pagination: modules.Pagination{
			Size: -1,
		},
	})
	require.Nil(t, err)
	require.Len(t, promotionFlowModel, 1)
}
