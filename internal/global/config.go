package global

import (
	"github.com/eden-framework/courier/transport_grpc"
	"github.com/eden-framework/courier/transport_http"
	"github.com/eden-framework/eden-framework/pkg/client/mysql"
	"github.com/eden-w2w/lib-modules/clients/gaode"
	"github.com/eden-w2w/lib-modules/databases"
	"github.com/eden-w2w/lib-modules/modules/id_generator"
	"github.com/eden-w2w/lib-modules/modules/settlement_flow"
	"github.com/eden-w2w/lib-modules/modules/wechat"
	"github.com/sirupsen/logrus"
	"time"
)

var Config = struct {
	LogLevel logrus.Level

	// db
	MasterDB *mysql.MySQL
	SlaveDB  *mysql.MySQL

	// administrator
	GRPCServer *transport_grpc.ServeGRPC
	HTTPServer *transport_http.ServeHTTP

	// id generation
	id_generator.SnowflakeConfig

	// wechat config
	Wechat wechat.Wechat

	// 订单超时时间
	OrderExpireIn time.Duration
	// 支付流水默认超时时间
	PaymentFlowExpireIn time.Duration
	// 结算配置
	settlement_flow.SettlementConfig

	ClientGaode *gaode.GaodeClient
}{
	LogLevel: logrus.DebugLevel,

	MasterDB: &mysql.MySQL{Database: databases.Config.DB},
	SlaveDB:  &mysql.MySQL{Database: databases.Config.DB},

	GRPCServer: &transport_grpc.ServeGRPC{
		Port: 8900,
	},
	HTTPServer: &transport_http.ServeHTTP{
		Port:     8800,
		WithCORS: true,
	},
	SnowflakeConfig: id_generator.SnowflakeConfig{
		Epoch:      1288351723598,
		BaseNodeID: 1,
		NodeCount:  100,
		NodeBits:   10,
		StepBits:   12,
	},
	OrderExpireIn:       30 * time.Minute,
	PaymentFlowExpireIn: 5 * time.Minute,
	SettlementConfig: settlement_flow.SettlementConfig{
		SettlementRules: []settlement_flow.SettlementRule{
			{
				MinSales:   0,
				MaxSales:   500000,
				Proportion: 0.1,
			},
			{
				MinSales:   500000,
				MaxSales:   5000000,
				Proportion: 0.15,
			},
			{
				MinSales:   5000000,
				MaxSales:   ^uint64(0),
				Proportion: 0.2,
			},
		},
		// 结算等待7天，可能涉及7天内退货
		SettlementDuration: 7 * 24 * time.Hour,
	},
	ClientGaode: &gaode.GaodeClient{},
}
