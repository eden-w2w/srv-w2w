package databases

import (
	"github.com/eden-framework/sqlx/datatypes"
	"github.com/eden-w2w/srv-w2w/internal/contants/enums"
)

//go:generate eden generate model StatementsFlow --database Config.DB --with-comments
//go:generate eden generate tag StatementsFlow --defaults=true
// @def primary ID
// @def unique_index U_statements_id StatementsID
// @def unique_index U_interval UserID Name
type StatementsFlow struct {
	datatypes.PrimaryID
	// 结算单ID
	StatementsID uint64 `json:"statementsID,string" db:"f_statements_id"`
	// 用户ID
	UserID uint64 `json:"userID,string" db:"f_user_id"`
	// 名称
	Name string `json:"name" db:"f_name"`
	// 销售总额
	TotalSales uint64 `json:"totalSales" db:"f_total_sales"`
	// 计算比例
	Proportion float64 `json:"proportion" db:"f_proportion"`
	// 结算金额
	Amount uint64 `json:"amount" db:"f_amount"`
	// 结算状态
	Status enums.StatementStatus `json:"status" db:"f_status"`

	datatypes.OperateTime
}
