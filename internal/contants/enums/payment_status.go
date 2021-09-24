package enums

//go:generate eden generate enum --type-name=PaymentStatus
// api:enum
type PaymentStatus uint8

// 支付状态
const (
	PAYMENT_STATUS_UNKNOWN  PaymentStatus = iota
	PAYMENT_STATUS__CREATED               // 未支付
	PAYMENT_STATUS__PROCESS               // 处理中
	PAYMENT_STATUS__SUCCESS               // 支付成功
	PAYMENT_STATUS__FAIL                  // 支付失败
)
