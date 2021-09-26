package promotion

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-w2w/srv-w2w/internal/contants/errors"
	"github.com/eden-w2w/srv-w2w/internal/modules"
	"github.com/eden-w2w/srv-w2w/internal/modules/promotion_flow"
	"github.com/eden-w2w/srv-w2w/internal/modules/settlement_flow"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
)

func init() {
	Router.Register(courier.NewRouter(GetMyPromotionSummary{}))
}

// GetMyPromotionSummary 获取我的推广概览
type GetMyPromotionSummary struct {
	httpx.MethodGet
}

func (req GetMyPromotionSummary) Path() string {
	return ""
}

func (req GetMyPromotionSummary) Output(ctx context.Context) (result interface{}, err error) {
	user := middleware.GetUserByContext(ctx)
	if user == nil {
		return nil, errors.Unauthorized
	}

	flows, err := promotion_flow.GetController().GetPromotionFlows(promotion_flow.GetPromotionFlowParams{
		UserID:          user.UserID,
		IsNotSettlement: true,
		Pagination: modules.Pagination{
			Size: -1,
		},
	})
	if err != nil {
		return nil, err
	}

	resp := &GetMyPromotionSummaryResponse{}
	resp.TotalSales, resp.ExpectedAmount = settlement_flow.GetController().GetPromotionSettlementAmount(flows)
	return resp, nil
}

type GetMyPromotionSummaryResponse struct {
	// 本期待结算金额
	TotalSales uint64 `json:"totalSales"`
	// 预期收入
	ExpectedAmount uint64 `json:"expectedAmount"`
}
