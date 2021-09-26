package test

import (
	"context"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
	"github.com/eden-w2w/srv-w2w/internal/routers/v0/promotion"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testGetMyPromotionSummary(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.AuthContextKey, promotionUserModel)

	request := promotion.GetMyPromotionSummary{}
	resp, err := request.Output(ctx)
	assert.Nil(t, err)

	response := resp.(*promotion.GetMyPromotionSummaryResponse)
	assert.Equal(t, orderModel.TotalPrice, response.TotalSales)

	var expectedAmount uint64 = 0
	for _, rule := range global.Config.SettlementRules {
		if orderModel.TotalPrice >= rule.MinSales && orderModel.TotalPrice < rule.MaxSales {
			expectedAmount = uint64(float64(orderModel.TotalPrice) * rule.Proportion)
		}
	}
	assert.Equal(t, expectedAmount, response.ExpectedAmount)
}
