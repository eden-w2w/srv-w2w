package wechat

import (
	"context"
	"fmt"
	errors "github.com/eden-w2w/lib-modules/constants/general_errors"
	"github.com/eden-w2w/lib-modules/databases"
	"github.com/eden-w2w/srv-w2w/internal/global"
	w "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/auth"
	programConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
	"github.com/sirupsen/logrus"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"net/http"
	"time"
)

var controller *Controller

func GetController() *Controller {
	if controller == nil {
		controller = newController(w.NewWechat(), global.Config.Wechat)
	}
	return controller
}

type Controller struct {
	wc        *w.Wechat
	program   *miniprogram.MiniProgram
	payClient *core.Client

	appID             string
	merchantID        string
	merchantKeySecret string
	defaultProdDesc   string
	notifyUrl         string
}

func newController(wc *w.Wechat, wechatConfig global.Wechat) *Controller {
	memory := cache.NewMemory()
	program := wc.GetMiniProgram(&programConfig.Config{
		AppID:     wechatConfig.AppID,
		AppSecret: wechatConfig.AppSecret,
		Cache:     memory,
	})

	var client *core.Client
	if global.Config.EnableWechatPay {
		mchPK, err := utils.LoadPrivateKey(wechatConfig.MerchantPK)
		if err != nil {
			logrus.Panicf("[wechat.newController] utils.LoadPrivateKey err: %v", err)
		}
		ctx := context.Background()
		opts := []core.ClientOption{
			option.WithWechatPayAutoAuthCipher(
				wechatConfig.MerchantID,
				wechatConfig.MerchantCertSerialNo,
				mchPK,
				wechatConfig.MerchantSecret),
		}
		client, err = core.NewClient(ctx, opts...)
		if err != nil {
			logrus.Panicf("[wechat.newController] core.NewClient err: %v", err)
		}
	}

	return &Controller{
		wc:        wc,
		program:   program,
		payClient: client,

		appID:             wechatConfig.AppID,
		merchantID:        wechatConfig.MerchantID,
		merchantKeySecret: wechatConfig.MerchantSecret,
		defaultProdDesc:   wechatConfig.ProductionDesc,
		notifyUrl:         wechatConfig.NotifyUrl,
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

func (c Controller) CreatePrePayment(ctx context.Context, order *databases.Order, flow *databases.PaymentFlow, payer *databases.User) (resp *jsapi.PrepayWithRequestPaymentResponse, err error) {
	if c.payClient == nil {
		return
	}
	service := jsapi.JsapiApiService{
		Client: c.payClient,
	}
	request := jsapi.PrepayRequest{
		Appid:         core.String(c.appID),
		Mchid:         core.String(c.merchantID),
		Description:   core.String(c.defaultProdDesc),
		OutTradeNo:    core.String(fmt.Sprintf("%d", flow.FlowID)),
		TimeExpire:    core.Time(time.Time(flow.ExpiredAt)),
		Attach:        nil,
		NotifyUrl:     core.String(c.notifyUrl),
		GoodsTag:      nil,
		LimitPay:      nil,
		SupportFapiao: nil,
		Amount: &jsapi.Amount{
			Total:    core.Int64(int64(flow.Amount)),
			Currency: nil,
		},
		Payer: &jsapi.Payer{
			Openid: core.String(payer.OpenID),
		},
		Detail:     nil,
		SceneInfo:  nil,
		SettleInfo: nil,
	}
	resp, _, err = service.PrepayWithRequestPayment(ctx, request)
	if err != nil {
		logrus.Errorf("[CreatePrePayment] service.PrepayWithRequestPayment err: %v, request: %+v", err, request)
		return nil, errors.BadGateway
	}
	return
}

func (c Controller) ParseWechatPaymentNotify(ctx context.Context, request *http.Request) (*notify.Request, *payments.Transaction, error) {
	certVisitor := downloader.MgrInstance().GetCertificateVisitor(c.merchantID)
	handler := notify.NewNotifyHandler(c.merchantKeySecret, verifiers.NewSHA256WithRSAVerifier(certVisitor))

	transaction := new(payments.Transaction)
	notifyReq, err := handler.ParseNotifyRequest(ctx, request, transaction)
	if err != nil {
		logrus.Errorf("[ParseWechatPaymentNotify] handler.ParseNotifyRequest err: %v", err)
		return nil, nil, errors.InternalError
	}

	return notifyReq, transaction, nil
}
