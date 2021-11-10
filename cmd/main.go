package main

import (
	"github.com/eden-framework/context"
	"github.com/eden-framework/eden-framework/pkg/application"
	"github.com/eden-framework/sqlx/migration"
	"github.com/eden-w2w/lib-modules/modules/events"
	"github.com/eden-w2w/lib-modules/modules/goods"
	"github.com/eden-w2w/lib-modules/modules/id_generator"
	"github.com/eden-w2w/lib-modules/modules/order"
	"github.com/eden-w2w/lib-modules/modules/payment_flow"
	"github.com/eden-w2w/lib-modules/modules/promotion_flow"
	"github.com/eden-w2w/lib-modules/modules/settlement_flow"
	"github.com/eden-w2w/lib-modules/modules/user"
	"github.com/eden-w2w/lib-modules/modules/wechat"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/eden-w2w/lib-modules/databases"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/eden-w2w/srv-w2w/internal/routers"
)

var cmdMigrationDryRun bool

func main() {
	app := application.NewApplication(
		runner, false,
		application.WithConfig(&global.Config),
		application.WithConfig(&databases.Config),
	)

	cmdMigrate := &cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			migrate(
				&migration.MigrationOpts{
					DryRun: cmdMigrationDryRun,
				},
			)
		},
	}
	cmdMigrate.Flags().BoolVarP(&cmdMigrationDryRun, "dry", "d", false, "migrate --dry")
	app.AddCommand(cmdMigrate)

	app.Start()
}

func runner(ctx *context.WaitStopContext) error {
	initModules()

	go global.Config.GRPCServer.Serve(ctx, routers.Router)
	return global.Config.HTTPServer.Serve(ctx, routers.Router)
}

func migrate(opts *migration.MigrationOpts) {
	if err := migration.Migrate(global.Config.MasterDB, opts); err != nil {
		panic(err)
	}
}

func initModules() {
	logrus.SetLevel(global.Config.LogLevel)
	id_generator.GetGenerator().Init(global.Config.SnowflakeConfig)
	user.GetController().Init(global.Config.MasterDB)
	wechat.GetController().Init(global.Config.Wechat)
	goods.GetController().Init(global.Config.MasterDB)
	order.GetController().Init(
		global.Config.MasterDB,
		global.Config.OrderExpireIn,
		events.NewOrderEvent(global.Config.Wechat),
	)
	payment_flow.GetController().Init(global.Config.MasterDB, global.Config.PaymentFlowExpireIn)
	promotion_flow.GetController().Init(global.Config.MasterDB)
	settlement_flow.GetController().Init(global.Config.MasterDB, global.Config.SettlementConfig)
}
