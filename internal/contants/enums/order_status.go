package enums

//go:generate eden generate enum --type-name=OrderStatus
// api:enum
type OrderStatus uint8

// 订单状态
const (
	ORDER_STATUS_UNKNOWN   OrderStatus = iota
	ORDER_STATUS__CREATED              // 待支付
	ORDER_STATUS__PAID                 // 已支付
	ORDER_STATUS__CONFIRM              // 待发货
	ORDER_STATUS__DISPATCH             // 已发货
	ORDER_STATUS__COMPLETE             // 已完成
	ORDER_STATUS__CANCEL               // 已取消
)
