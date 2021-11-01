module github.com/eden-w2w/srv-w2w

go 1.16

replace k8s.io/client-go => k8s.io/client-go v0.18.8

require (
	github.com/eden-framework/context v0.0.2
	github.com/eden-framework/courier v1.0.5
	github.com/eden-framework/eden-framework v1.2.6-0.20211020014935-eab59ae7d198
	github.com/eden-framework/sqlx v0.0.1
	github.com/eden-w2w/lib-modules v0.0.9-0.20211031065607-e0a592e8aa80
	github.com/mozillazg/go-pinyin v0.18.0 // indirect
	github.com/silenceper/wechat/v2 v2.0.9
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.7.0
	github.com/wechatpay-apiv3/wechatpay-go v0.2.7
)
