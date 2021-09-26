package global

import (
	"fmt"
	"github.com/eden-framework/courier/transport_grpc"
	"github.com/eden-framework/courier/transport_http"
	"github.com/eden-framework/eden-framework/pkg/client/mysql"
	"github.com/eden-w2w/srv-w2w/internal/contants/enums"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/eden-w2w/srv-w2w/internal/databases"
)

type SnowflakeConfig struct {
	Epoch      int64
	BaseNodeID int64
	NodeCount  int64
	NodeBits   uint8
	StepBits   uint8
}

type Wechat struct {
	// 小程序AppID
	AppID string
	// 小程序AppSecret
	AppSecret string
	// 微信商户ID
	MerchantID string
	// 微信商户证书序列号
	MerchantCertSerialNo string
	// 微信商户证书私钥
	MerchantPK string
	// 微信商户APIv3密钥
	MerchantSecret string
	// 微信支付商品描述
	ProductionDesc string
	// 微信支付回调地址
	NotifyUrl string
	// 启用微信支付
	EnableWechatPay bool
}

type StatementConfig struct {
	// 结算周期
	StatementType enums.StatementType
	// 结算节点 周：0-6，月：1-31
	StatementDate uint8
}

func (c StatementConfig) ToStatementCronRule() string {
	if c.StatementType == enums.STATEMENT_TYPE__WEEK {
		return fmt.Sprintf("0 0 0 * * %d", c.StatementDate)
	} else if c.StatementType == enums.STATEMENT_TYPE__MONTH {
		return fmt.Sprintf("0 0 0 %d * *", c.StatementDate)
	}
	return ""
}

var Config = struct {
	LogLevel logrus.Level

	// db
	MasterDB *mysql.MySQL
	SlaveDB  *mysql.MySQL

	// administrator
	GRPCServer *transport_grpc.ServeGRPC
	HTTPServer *transport_http.ServeHTTP

	// id generation
	SnowflakeConfig

	// wechat config
	Wechat

	// 订单超时时间
	OrderExpireIn time.Duration
	// 支付流水默认超时时间
	PaymentFlowExpireIn time.Duration

	StatementConfig
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
	SnowflakeConfig: SnowflakeConfig{
		Epoch:      1288351723598,
		BaseNodeID: 1,
		NodeCount:  100,
		NodeBits:   10,
		StepBits:   12,
	},
	OrderExpireIn:       30 * time.Minute,
	PaymentFlowExpireIn: 5 * time.Minute,
	StatementConfig: StatementConfig{
		StatementType: enums.STATEMENT_TYPE__WEEK,
		StatementDate: 1,
	},
}
