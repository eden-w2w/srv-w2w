package promotion_flow

type CreatePromotionFlowParams struct {
	// 获得奖励的用户ID
	UserID uint64 `json:"userID,string"`
	// 获得奖励的用户昵称
	UserNickName string `json:"userNickName"`
	// 奖励来源用户ID
	RefererID uint64 `json:"refererID,string"`
	// 奖励来源的用户昵称
	RefererNickName string `json:"refererNickName"`
	// 奖励金额
	Amount uint64 `json:"amount"`
	// 奖励比例
	Proportion float64 `json:"proportion"`
	// 关联的支付流水
	PaymentFlowID uint64 `json:"paymentFlowID"`
	// 关联的结算单
	StatementID uint64 `json:"statementID,string"`
}
