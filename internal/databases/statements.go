package databases

import (
	"github.com/eden-framework/sqlx/datatypes"
	"github.com/eden-w2w/srv-w2w/internal/contants/enums"
)

//go:generate eden generate model Statements --database Config.DB --with-comments
//go:generate eden generate tag Statements --defaults=true
// @def primary ID
// @def unique_index U_statements_id StatementsID
// @def unique_index U_interval UserID StartTime EndTime
type Statements struct {
	datatypes.PrimaryID
	// 结算单ID
	StatementsID uint64 `json:"statementsID,string" db:"f_statements_id"`
	// 用户ID
	UserID uint64 `json:"userID,string" db:"f_user_id"`
	// 结算周期开始时间
	StartTime datatypes.MySQLTimestamp `json:"startTime" db:"f_start_time"`
	// 结算周期结束时间
	EndTime datatypes.MySQLTimestamp `json:"endTime" db:"f_end_time"`
	// 结算金额
	Amount uint64 `json:"amount" db:"f_amount"`
	// 结算状态
	Status enums.StatementStatus `json:"status" db:"f_status"`

	datatypes.OperateTime
}
