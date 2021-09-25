package statement

import "github.com/eden-framework/sqlx/datatypes"

type CreateStatementParams struct {
	// 用户ID
	UserID uint64 `in:"body" json:"userID,string" default:""`
	// 结算周期开始时间
	StartTime datatypes.MySQLTimestamp `in:"body" json:"startTime"`
	// 结算周期结束时间
	EndTime datatypes.MySQLTimestamp `in:"body" json:"endTime"`
	// 结算金额
	Amount uint64 `in:"body" json:"amount"`
}
