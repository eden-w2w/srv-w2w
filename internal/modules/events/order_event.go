package events

import (
	"github.com/eden-framework/sqlx"
	"github.com/eden-w2w/srv-w2w/internal/databases"
	"github.com/eden-w2w/srv-w2w/internal/modules/payment_flow"
	"github.com/eden-w2w/srv-w2w/internal/modules/promotion_flow"
	"github.com/eden-w2w/srv-w2w/internal/modules/user"
)

type OrderEvent struct {
}

func (o *OrderEvent) OnOrderCreateEvent(db sqlx.DBExecutor, order *databases.Order) error {
	return nil
}

func (o *OrderEvent) OnOrderPaidEvent(db sqlx.DBExecutor, order *databases.Order, payment *databases.PaymentFlow) error {
	return nil
}

func (o *OrderEvent) OnOrderCompleteEvent(db sqlx.DBExecutor, order *databases.Order) error {
	// 获取支付流水
	flow, err := payment_flow.GetController().GetFlowByOrderAndUserID(order.OrderID, order.UserID, db)
	if err != nil {
		return err
	}

	// 获取订单创建者
	orderUser, err := user.GetController().GetUserByUserID(order.UserID, db, true)
	if err != nil {
		return err
	}

	// 如果创建者没有推荐者信息则无需计算提成
	if orderUser.RefererID == 0 {
		return nil
	}
	// 获取推荐者
	refererUser, err := user.GetController().GetUserByUserID(orderUser.RefererID, db, true)
	if err != nil {
		return err
	}

	// 计算提成并创建提成流水
	proCtrl := promotion_flow.GetController()
	proportion := proCtrl.GetProportion()
	proAmount := uint64(float64(flow.Amount) * proportion) // 提成金额，单位分，向下取整
	_, err = proCtrl.CreatePromotionFlow(promotion_flow.CreatePromotionFlowParams{
		UserID:          refererUser.UserID,
		UserNickName:    refererUser.NickName,
		RefererID:       orderUser.UserID,
		RefererNickName: orderUser.NickName,
		Amount:          proAmount,
		Proportion:      proportion,
		PaymentFlowID:   flow.FlowID,
	}, db)
	return err
}

func NewOrderEvent() *OrderEvent {
	return &OrderEvent{}
}
