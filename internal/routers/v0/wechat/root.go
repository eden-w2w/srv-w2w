package wechat

import "github.com/eden-framework/courier"

var Router = courier.NewRouter(WechatRouter{})

type WechatRouter struct {
	courier.EmptyOperator
}

func (WechatRouter) Path() string {
	return "/wechat"
}
