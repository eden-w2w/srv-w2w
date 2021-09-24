package main

import (
	"github.com/eden-framework/context"
	"github.com/eden-framework/eden-framework/pkg/application"
	"github.com/eden-framework/sqlx/migration"
	"github.com/eden-w2w/srv-w2w/internal/modules/events"
	"github.com/eden-w2w/srv-w2w/internal/modules/goods"
	"github.com/eden-w2w/srv-w2w/internal/modules/order"
	"github.com/eden-w2w/srv-w2w/internal/modules/payment_flow"
	"github.com/eden-w2w/srv-w2w/internal/modules/promotion_flow"
	"github.com/eden-w2w/srv-w2w/internal/modules/user"
	"github.com/eden-w2w/srv-w2w/internal/modules/wechat"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/eden-w2w/srv-w2w/internal"
	"github.com/eden-w2w/srv-w2w/internal/databases"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/eden-w2w/srv-w2w/internal/routers"
)

var cmdMigrationDryRun bool

func main() {
	app := application.NewApplication(runner,
		application.WithConfig(&global.Config),
		application.WithConfig(&databases.Config))

	cmdMigrate := &cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			migrate(&migration.MigrationOpts{
				DryRun: cmdMigrationDryRun,
			})
		},
	}
	cmdMigrate.Flags().BoolVarP(&cmdMigrationDryRun, "dry", "d", false, "migrate --dry")
	app.AddCommand(cmdMigrate)

	app.Start()
}

func runner(ctx *context.WaitStopContext) error {
	logrus.SetLevel(global.Config.LogLevel)
	internal.GetGenerator()
	user.GetController()
	wechat.GetController()
	goods.GetController()
	order.GetController().WithEventHandler(events.NewOrderEvent())
	payment_flow.GetController()
	promotion_flow.GetController()

	go global.Config.GRPCServer.Serve(ctx, routers.Router)
	return global.Config.HTTPServer.Serve(ctx, routers.Router)
}

func migrate(opts *migration.MigrationOpts) {
	if err := migration.Migrate(global.Config.MasterDB, opts); err != nil {
		panic(err)
	}
}
