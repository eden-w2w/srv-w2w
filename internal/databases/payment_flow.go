package databases

import (
	"github.com/eden-framework/sqlx/datatypes"
	"github.com/eden-w2w/srv-w2w/internal/contants/enums"
)

//go:generate eden generate model PaymentFlow --database Config.DB --with-comments
//go:generate eden generate tag PaymentFlow --defaults=true
// @def primary ID
// @def unique_index U_flow_id FlowID
type PaymentFlow struct {
	datatypes.PrimaryID
	// 流水ID
	FlowID uint64 `json:"flowID,string" db:"f_flow_id"`
	// 用户ID
	UserID uint64 `json:"userID,string" db:"f_user_id"`
	// 关联订单号
	OrderID uint64 `json:"orderID,string" db:"f_order_id"`
	// 支付金额
	Amount uint64 `json:"amount" db:"f_amount"`
	// 支付方式
	PaymentMethod enums.PaymentMethod `json:"paymentMethod" db:"f_payment_method"`
	// 支付系统流水号
	RemoteFlowID string `json:"remoteFlowID" db:"f_remote_flow_id,size=255"`

	datatypes.OperateTime
}
