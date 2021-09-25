package global

import (
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

	// 全局默认提成比例
	GlobalProportion float64
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
	GlobalProportion:    0.005,
	PaymentFlowExpireIn: 30 * time.Minute,
	Wechat: Wechat{
		AppID:                "wx6f07849ea76cb144",
		AppSecret:            "d0ba9e67879050de041671562e8e57a7",
		MerchantID:           "asdfasdfasdfsdaf",
		MerchantCertSerialNo: "asdfasdfasdfasdf",
		MerchantPK:           "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDcfx2unKTI5wo5\nX2AW8MSAppgvI7TbY1tpXiYcQ8LlIPMr5jTT9PImA9xXL+JNm4n+qTXNxee049UZ\nRv6P5kgimvw+NB3UkQ361KQP41ZInnuLcCBcVRXASie+MRxSucno4c7Y/PtFtbZp\nMll+v4uziHP3BRuY26CPzrLMPy6rkKlfxN3hf1MpCCgBcKyd2niF7BBhYavhmQA2\nryrP3V+8eE5XVPLeUZ/wbV7QX1w4Hyf99lM4CWzlm631OhJnCkM7VebCEJE/pI5b\nOZT9VduP2H+ZJ2HU7FY4ErL6gaoc2+eaUy21VooNUrlNIX1NxeH+KweeLAZe9cMF\nKp6pK5bDAgMBAAECggEAB4LbB341jNR2ADNUCiQEhZ9nMr/kIs8VJ5c59B+IzZZc\nQkQSy5+d+haKlHvAmUUxo7eu6gYSWRKjyTsmY8+D6I3gBtxuobWZRNCmkBX98/ZW\nHg9hQtPxLD46cVup69WBX8oFalXf12WTw6yf6NXsk96TmLgXUxM0OoHfjF6gn956\n9cF8YoBfMNMpEBmwYwncEuN5HjmBLBDJgA7emhs+jcGgDnkdggOFNHhHBWBlJUTC\nqDdZ6m78Ek/sbUdBtSldKHxWB7j0LIemAdlhfKJzxffoLLq2YBCxS6DbiRpjA8fJ\nysbWixkAI4kmgJrGzkV1Jqv+Xtsx48MGmz0v1R5VaQKBgQD9wr+7huixubNA9MEd\nST4g8NY3HVUIooctIB1Cfa1b76+AGrNSq+HJAx1Ws6zThYKZZNFYurj1j3ibhgIG\nelneYA3Ig8Kv8Rid3tfmRfV+bwH8uqNy8aAviEzdvPbp8wCF8KJZYyHnZVitz88E\nFwXJlXYPhpKR+ZV1fTkBHxceqwKBgQDecTjy4c/UJ1YP5j4YgbR3mqQ11sFPwkGi\nXw5c8V3dRVbJwOISggOuVcdYtQuv5V0YjGpH+Tnu2HWqASwfXtWf/GwCByZjZsMm\nCXAmUKOveF4XdydTToF9DatwrfwH3sVbSo6tr0/0/BxxABIncuv3I9TI+uytOENr\n1K531duISQKBgQCGea/1hrbFiC1QHORBytCb0EUVC/xGCSstZLlcxREbiVctwfiJ\ntQB/76Cak7jglv3woBa4uKrPpuo5MLjeCfSZ/bkQK5L+ffuXncI4C5bfG6Cn95gq\nj5Vd2xMw5rTKobYDRNQkHn6XC4QFB/0io7izbPsVmzENHoBvtJ6C06EC5QKBgQDE\nlQ8dJw8CkRjvyCzMf6Q7p0DEC9yfXi+fDZ6l6PFpIWRHvNyOOrSCF3CrJqUDeBJE\nJaOvo8PdHxMtgQe7WZY1Y2EYkbUNV73TGbgxREbERf3xgu4QI2swVypR1JDNa8TC\nnytkt3BUM5H3E3b3wgpjSRk93PPZScXbwAPICAmRcQKBgB/krZFrIqzRcKzJEsz3\nynWi8NIFeSWh5K5yU46w/rsdTKTK6ILkLLZkHgQCNEWuRZKUK6z0GM5Fj3C1pyvh\nAEOnhMSNP+J8IrUtVnRo65dubqadvwGOmbEOfraXdI7U+npS1UWZlHqz+xnkZrTf\nNNfQ5DKd/ueHFVXQ+Cjt1Tc2\n-----END PRIVATE KEY-----\n",
		MerchantSecret:       "7e577d434ffb7118b51058fd3d9056f7",
		ProductionDesc:       "比翼婚宴酒",
		NotifyUrl:            "https://api.w2wing.com/w2w/v0/wechat/notify",
		EnableWechatPay:      false,
	},
	StatementConfig: StatementConfig{
		StatementType: enums.STATEMENT_TYPE__WEEK,
		StatementDate: 1,
	},
}
