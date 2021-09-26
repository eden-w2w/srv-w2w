package statement_flow

type CreateStatementParams struct {
	// 用户ID
	UserID uint64 `in:"body" json:"userID,string" default:""`
	// 名称
	Name string `in:"body" json:"name"`
	// 销售总额
	TotalSales uint64 `in:"body" json:"totalSales"`
	// 计算比例
	Proportion float64 `in:"body" json:"proportion"`
	// 结算金额
	Amount uint64 `in:"body" json:"amount"`
}
