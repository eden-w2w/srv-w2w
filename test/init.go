package test

import (
	"github.com/eden-framework/context"
	"github.com/eden-framework/eden-framework/pkg/application"
	"github.com/eden-w2w/lib-modules/databases"
	"github.com/eden-w2w/lib-modules/modules/events"
	"github.com/eden-w2w/lib-modules/modules/goods"
	"github.com/eden-w2w/lib-modules/modules/id_generator"
	"github.com/eden-w2w/lib-modules/modules/order"
	"github.com/eden-w2w/lib-modules/modules/payment_flow"
	"github.com/eden-w2w/lib-modules/modules/promotion_flow"
	"github.com/eden-w2w/lib-modules/modules/settlement_flow"
	"github.com/eden-w2w/lib-modules/modules/user"
	"github.com/eden-w2w/srv-w2w/internal/global"
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
	id_generator.GetGenerator().Init(global.Config.SnowflakeConfig)
	user.GetController().Init(global.Config.MasterDB)
	wechat.GetController()
	goods.GetController().Init(global.Config.MasterDB)
	order.GetController().Init(global.Config.MasterDB, global.Config.OrderExpireIn, events.NewOrderEvent())
	payment_flow.GetController().Init(global.Config.MasterDB, global.Config.PaymentFlowExpireIn)
	promotion_flow.GetController().Init(global.Config.MasterDB)
	settlement_flow.GetController().Init(global.Config.MasterDB, global.Config.SettlementConfig)

	return nil
}
