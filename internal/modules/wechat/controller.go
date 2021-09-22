package wechat

import (
	"github.com/eden-w2w/srv-w2w/internal/contants/errors"
	"github.com/eden-w2w/srv-w2w/internal/global"
	w "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/auth"
	"github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
	"github.com/sirupsen/logrus"
)

var controller *Controller

func GetController() *Controller {
	if controller == nil {
		controller = NewController(w.NewWechat(), global.Config.WechatAppID, global.Config.WechatAppSecret)
	}
	return controller
}

type Controller struct {
	wc      *w.Wechat
	program *miniprogram.MiniProgram
}

func NewController(wc *w.Wechat, appID, appSecret string) *Controller {
	memory := cache.NewMemory()
	cfg := &config.Config{
		AppID:     appID,
		AppSecret: appSecret,
		Cache:     memory,
	}
	program := wc.GetMiniProgram(cfg)

	return &Controller{
		wc:      wc,
		program: program,
	}
}

func (c Controller) Code2Session(code string) (*auth.ResCode2Session, error) {
	resp, err := c.program.GetAuth().Code2Session(code)
	if err != nil {
		logrus.Errorf("[Code2Session] c.program.GetAuth().Code2Session(code) err: %v, code: %s", err, code)
		return nil, errors.BadGateway
	}
	return &resp, nil
}

func (c Controller) GetUnlimitedQrCode(params qrcode.QRCoder) (buffer []byte, err error) {
	buffer, err = c.program.GetQRCode().GetWXACodeUnlimit(params)
	if err != nil {
		logrus.Errorf("[GetUnlimitedQrCode] c.program.GetQRCode().GetWXACodeUnlimit(params) err: %v, params: %+v", err, params)
		return nil, errors.BadGateway
	}

	return
}
