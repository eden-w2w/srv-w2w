package orders

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-framework/sqlx"
	"github.com/eden-w2w/lib-modules/constants/enums"
	"github.com/eden-w2w/lib-modules/modules/goods"
	"github.com/eden-w2w/lib-modules/modules/order"
	"github.com/eden-w2w/srv-w2w/internal/constants/errors"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
)

func init() {
	Router.Register(courier.NewRouter(CompleteOrder{}))
}

// CompleteOrder 确认收货
type CompleteOrder struct {
	httpx.MethodPut

	// 订单ID
	OrderID uint64 `in:"path" name:"orderID,string"`
}

func (req CompleteOrder) Path() string {
	return "/:orderID/complete"
}

func (req CompleteOrder) Output(ctx context.Context) (result interface{}, err error) {
	user := middleware.GetUserByContext(ctx)
	if user == nil {
		return nil, errors.Unauthorized
	}

	tx := sqlx.NewTasks(global.Config.MasterDB)
	tx = tx.With(
		func(db sqlx.DBExecutor) error {
			orderModel, logistics, err := order.GetController().GetOrder(req.OrderID, user.UserID, db, true)
			if err != nil {
				return err
			}
			orderGoods, err := order.GetController().GetOrderGoods(req.OrderID, db)
			if err != nil {
				return err
			}
			return order.GetController().UpdateOrder(
				orderModel,
				logistics,
				orderGoods,
				order.UpdateOrderParams{
					Status: enums.ORDER_STATUS__COMPLETE,
				},
				goods.GetController().LockInventory,
				goods.GetController().UnlockInventory,
				db,
			)
		},
	)

	err = tx.Do()
	return
}
