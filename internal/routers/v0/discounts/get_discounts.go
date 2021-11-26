package discounts

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-w2w/lib-modules/constants/enums"
	"github.com/eden-w2w/lib-modules/modules"
	"github.com/eden-w2w/lib-modules/modules/discounts"
	"time"
)

func init() {
	Router.Register(courier.NewRouter(GetDiscounts{}))
}

// GetDiscounts 获取优惠列表
type GetDiscounts struct {
	httpx.MethodGet
	// 商品ID
	GoodsID uint64 `in:"query" name:"goodsID,string" default:""`
}

func (req GetDiscounts) Path() string {
	return ""
}

func (req GetDiscounts) Output(ctx context.Context) (result interface{}, err error) {
	list, _, err := discounts.GetController().GetDiscounts(
		discounts.GetDiscountsParams{
			Status: enums.DISCOUNT_STATUS__PROCESS,
			Pagination: modules.Pagination{
				Size: -1,
			},
		}, false,
	)

	// 过滤无效优惠
	for i := 0; i < len(list); i++ {
		func() {
			_ = discounts.GetController().RLock(list[i].DiscountID)
			defer discounts.GetController().RUnlock(list[i].DiscountID)

			current := time.Now()
			if list[i].Limit > 0 && list[i].Times >= uint64(list[i].Limit) {
				list = append(list[:i], list[i+1:]...)
			} else if current.Before(time.Time(list[i].ValidityStart)) {
				list = append(list[:i], list[i+1:]...)
			} else if current.After(time.Time(list[i].ValidityEnd)) {
				list = append(list[:i], list[i+1:]...)
			}
		}()
	}
	return list, err
}
