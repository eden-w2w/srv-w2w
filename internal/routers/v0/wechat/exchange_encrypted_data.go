package wechat

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	errors "github.com/eden-w2w/lib-modules/constants/general_errors"
	"github.com/eden-w2w/lib-modules/modules/user"
	"github.com/eden-w2w/srv-w2w/internal/modules/wechat"
	"github.com/eden-w2w/srv-w2w/internal/routers/middleware"
)

func init() {
	Router.Register(courier.NewRouter(middleware.Authorization{}, ExchangeEncryptedData{}))
}

// ExchangeEncryptedData 获取加密解密数据
type ExchangeEncryptedData struct {
	httpx.MethodPost

	Data wechat.WechatUserInfo `in:"body"`
}

func (req ExchangeEncryptedData) Path() string {
	return "/exchange_user_data"
}

func (req ExchangeEncryptedData) Output(ctx context.Context) (result interface{}, err error) {
	u := middleware.GetUserByContext(ctx)
	if u == nil {
		return nil, errors.Unauthorized
	}

	plain, err := wechat.GetController().ExchangeEncryptedData(u.SessionKey, req.Data)
	if err != nil {
		return
	}

	err = user.GetController().UpdateUserInfo(u.UserID, user.UpdateUserInfoParams{
		Mobile:    plain.PurePhoneNumber,
		NickName:  plain.NickName,
		AvatarUrl: plain.AvatarURL,
	})
	if err != nil {
		return
	}
	return plain, nil
}
