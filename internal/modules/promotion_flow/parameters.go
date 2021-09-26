package promotion_flow

import (
	"github.com/eden-framework/sqlx/builder"
	"github.com/eden-w2w/srv-w2w/internal/databases"
	"github.com/eden-w2w/srv-w2w/internal/modules"
)

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
	// 关联的支付流水
	PaymentFlowID uint64 `json:"paymentFlowID"`
}

type GetPromotionFlowParams struct {
	// 获得奖励的用户ID
	UserID uint64 `name:"userID,string" in:"query"`
	// 奖励来源用户ID
	RefererID uint64 `name:"refererID,string" in:"query"`
	// 关联的支付流水
	PaymentFlowID uint64 `name:"paymentFlowID,string" in:"query"`
	// 关联的结算单ID
	StatementsID uint64 `name:"statementsID,string" in:"query"`
	modules.Pagination
}

func (p GetPromotionFlowParams) Conditions() builder.SqlCondition {
	var condition builder.SqlCondition
	var model = databases.PromotionFlow{}
	if p.UserID != 0 {
		condition = builder.And(condition, model.FieldUserID().Eq(p.UserID))
	}
	if p.RefererID != 0 {
		condition = builder.And(condition, model.FieldRefererID().Eq(p.RefererID))
	}
	if p.PaymentFlowID != 0 {
		condition = builder.And(condition, model.FieldPaymentFlowID().Eq(p.PaymentFlowID))
	}
	if p.StatementsID != 0 {
		condition = builder.And(condition, model.FieldStatementsID().Eq(p.StatementsID))
	}
	return condition
}

func (p GetPromotionFlowParams) Additions() []builder.Addition {
	var additions = make([]builder.Addition, 0)

	if p.Size != 0 {
		limit := builder.Limit(int64(p.Size))
		if p.Offset != 0 {
			limit.Offset(int64(p.Offset))
		}
		additions = append(additions, limit)
	}

	additions = append(additions, builder.OrderBy(builder.DescOrder((&databases.Order{}).FieldCreatedAt())))

	return additions
}
