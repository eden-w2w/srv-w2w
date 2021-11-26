package v0

import (
	"github.com/eden-framework/courier"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
	"github.com/eden-w2w/srv-w2w/internal/routers/v0/discounts"
	"github.com/eden-w2w/srv-w2w/internal/routers/v0/goods"
	"github.com/eden-w2w/srv-w2w/internal/routers/v0/orders"
	"github.com/eden-w2w/srv-w2w/internal/routers/v0/payment"
	"github.com/eden-w2w/srv-w2w/internal/routers/v0/promotion"
	"github.com/eden-w2w/srv-w2w/internal/routers/v0/settings"
	"github.com/eden-w2w/srv-w2w/internal/routers/v0/users"
	"github.com/eden-w2w/srv-w2w/internal/routers/v0/wechat"
)

var Router = courier.NewRouter(V0Router{})
var AuthRouter = courier.NewRouter(middleware.Authorization{})

type V0Router struct {
	courier.EmptyOperator
}

func (V0Router) Path() string {
	return "/v0"
}

func init() {
	Router.Register(wechat.Router)
	Router.Register(AuthRouter)

	Router.Register(goods.Router)
	AuthRouter.Register(orders.Router)
	AuthRouter.Register(users.Router)
	AuthRouter.Register(payment.Router)
	AuthRouter.Register(promotion.Router)
	AuthRouter.Register(settings.Router)
	AuthRouter.Register(discounts.Router)
}
