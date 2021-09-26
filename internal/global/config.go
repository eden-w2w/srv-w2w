package global

import (
	"errors"
	"fmt"
	"github.com/eden-framework/courier/transport_grpc"
	"github.com/eden-framework/courier/transport_http"
	"github.com/eden-framework/eden-framework/pkg/client/mysql"
	"github.com/eden-w2w/srv-w2w/internal/contants/enums"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
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

type SettlementConfig struct {
	// 结算周期
	SettlementType enums.SettlementType
	// 结算节点 周：0-6，月：1-31
	SettlementDate uint8
	// 提成比例规则
	SettlementRules []SettlementRule
}

func (c SettlementConfig) ToSettlementCronRule() string {
	if c.SettlementType == enums.SETTLEMENT_TYPE__WEEK {
		return fmt.Sprintf("0 0 0 * * %d", c.SettlementDate)
	} else if c.SettlementType == enums.SETTLEMENT_TYPE__MONTH {
		return fmt.Sprintf("0 0 0 %d * *", c.SettlementDate)
	}
	return ""
}

type SettlementRule struct {
	// 最小销售量（闭区间）
	MinSales uint64
	// 最大销售量（开区间）
	MaxSales uint64
	// 计提比例
	Proportion float64
}

func (s SettlementRule) String() string {
	str, _ := s.MarshalText()
	return string(str)
}

func (s *SettlementRule) UnmarshalText(text []byte) (err error) {
	strList := strings.Split(string(text), "|")
	if len(strList) != 3 {
		return errors.New("SettlementRule not support more than 3 args")
	}
	s.MinSales, err = strconv.ParseUint(strList[0], 10, 64)
	if err != nil {
		return
	}
	s.MaxSales, err = strconv.ParseUint(strList[1], 10, 64)
	if err != nil {
		return
	}
	s.Proportion, err = strconv.ParseFloat(strList[1], 64)
	return
}

func (s SettlementRule) MarshalText() (text []byte, err error) {
	str := fmt.Sprintf("%d|%d|%f", s.MinSales, s.MaxSales, s.Proportion)
	return []byte(str), nil
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

	SettlementConfig
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
	SettlementConfig: SettlementConfig{
		SettlementType: enums.SETTLEMENT_TYPE__WEEK,
		SettlementDate: 1,
		SettlementRules: []SettlementRule{
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
	},
}
