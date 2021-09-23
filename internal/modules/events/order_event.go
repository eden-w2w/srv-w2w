package events

import (
	"github.com/eden-framework/sqlx"
	"github.com/eden-w2w/srv-w2w/internal/databases"
)

var orderEvent *OrderEvent

type OrderEvent struct {
}

func (o *OrderEvent) OnOrderCreateEvent(db sqlx.DBExecutor, order *databases.Order) error {
	return nil
}

func (o *OrderEvent) OnOrderPaidEvent(db sqlx.DBExecutor, order *databases.Order) error {
	return nil
}

func (o *OrderEvent) OnOrderCompleteEvent(db sqlx.DBExecutor, order *databases.Order) error {
	return nil
}

func NewOrderEvent() *OrderEvent {
	return &OrderEvent{}
}

func GetOrderEvent() *OrderEvent {
	if orderEvent == nil {
		orderEvent = NewOrderEvent()
	}
	return orderEvent
}
