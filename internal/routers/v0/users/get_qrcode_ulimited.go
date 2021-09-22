package users

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-w2w/srv-w2w/internal/modules/wechat"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
)

func init() {
	Router.Register(courier.NewRouter(GetQrCodeUnlimited{}))
}

// GetQrCodeUnlimited 获取小程序码
type GetQrCodeUnlimited struct {
	httpx.MethodGet
	// page 必须是已经发布的小程序存在的页面,根路径前不要填加 /,不能携带参数（参数请放在scene字段里），如果不填写这个字段，默认跳主页面
	Page string `in:"query" default:"" name:"page,omitempty"`
	// path 扫码进入的小程序页面路径
	PagePath string `in:"query" default:"" name:"path,omitempty"`
	// scene 最大32个可见字符，只支持数字，大小写英文以及部分特殊字符：!#$&'()*+,/:;=?@-._~，其它字符请自行编码为合法字符（因不支持%，中文无法使用 urlencode 处理，请使用其他编码方式）
	Scene string `in:"query" json:"scene"`
	// width 图片宽度
	Width int `in:"query" default:"" name:"width,omitempty"`
	// autoColor 自动配置线条颜色，如果颜色依然是黑色，则说明不建议配置主色调
	AutoColor bool `in:"query" default:"" name:"autoColor,omitempty"`
	// lineColor AutoColor 为 false 时生效，使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"},十进制表示
	LineColorR string `in:"query" default:"" name:"lineColorR,omitempty"`
	// lineColor AutoColor 为 false 时生效，使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"},十进制表示
	LineColorG string `in:"query" default:"" name:"lineColorG,omitempty"`
	// lineColor AutoColor 为 false 时生效，使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"},十进制表示
	LineColorB string `in:"query" default:"" name:"lineColorB,omitempty"`
	// isHyaline 是否需要透明底色
	IsHyaline bool `in:"query" default:"" name:"isHyaline,omitempty"`
}

func (req GetQrCodeUnlimited) Path() string {
	return "/getQrCodeUnlimited"
}

func (req GetQrCodeUnlimited) Output(ctx context.Context) (result interface{}, err error) {
	return wechat.GetController().GetUnlimitedQrCode(qrcode.QRCoder{
		Page:      req.Page,
		Path:      req.PagePath,
		Width:     req.Width,
		Scene:     req.Scene,
		AutoColor: req.AutoColor,
		LineColor: qrcode.Color{
			R: req.LineColorR,
			G: req.LineColorG,
			B: req.LineColorB,
		},
		IsHyaline: req.IsHyaline,
	})
}
