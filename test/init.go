package test

import (
	"github.com/eden-framework/context"
	"github.com/eden-framework/eden-framework/pkg/application"
	"github.com/eden-w2w/srv-w2w/internal/databases"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/eden-w2w/srv-w2w/internal/modules/events"
	"github.com/eden-w2w/srv-w2w/internal/modules/goods"
	"github.com/eden-w2w/srv-w2w/internal/modules/id_generator"
	"github.com/eden-w2w/srv-w2w/internal/modules/order"
	"github.com/eden-w2w/srv-w2w/internal/modules/payment_flow"
	"github.com/eden-w2w/srv-w2w/internal/modules/promotion_flow"
	"github.com/eden-w2w/srv-w2w/internal/modules/settlement_flow"
	"github.com/eden-w2w/srv-w2w/internal/modules/user"
	"github.com/eden-w2w/srv-w2w/internal/modules/wechat"
	"github.com/sirupsen/logrus"
)

func init() {
	app := application.NewApplication(runner, true,
		application.WithConfig(&global.Config),
		application.WithConfig(&databases.Config))

	app.Start()
}

func runner(ctx *context.WaitStopContext) error {
	logrus.SetLevel(global.Config.LogLevel)
	id_generator.GetGenerator()
	user.GetController()
	wechat.GetController()
	goods.GetController()
	order.GetController().WithEventHandler(events.NewOrderEvent())
	payment_flow.GetController()
	promotion_flow.GetController()
	settlement_flow.GetController()

	return nil
}
