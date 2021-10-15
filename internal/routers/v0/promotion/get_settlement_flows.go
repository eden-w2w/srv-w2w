package promotion

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	errors "github.com/eden-w2w/lib-modules/constants/general_errors"
	"github.com/eden-w2w/lib-modules/modules"
	"github.com/eden-w2w/lib-modules/modules/settlement_flow"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
)

func init() {
	Router.Register(courier.NewRouter(GetSettlementFlows{}))
}

// GetSettlementFlows 获取结算流水单
type GetSettlementFlows struct {
	httpx.MethodGet
	modules.Pagination
}

func (req GetSettlementFlows) Path() string {
	return "/settlement_flows"
}

func (req GetSettlementFlows) Output(ctx context.Context) (result interface{}, err error) {
	user := middleware.GetUserByContext(ctx)
	if user == nil {
		return nil, errors.Unauthorized
	}

	list, _, err := settlement_flow.GetController().GetSettlementFlows(settlement_flow.GetSettlementFlowsParams{
		UserID:     user.UserID,
		Pagination: req.Pagination,
	}, false)
	return list, err
}
