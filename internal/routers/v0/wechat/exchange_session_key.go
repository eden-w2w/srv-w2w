package wechat

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	errors "github.com/eden-w2w/lib-modules/constants/general_errors"
	"github.com/eden-w2w/lib-modules/modules/user"
	"github.com/eden-w2w/lib-modules/modules/wechat"
)

func init() {
	Router.Register(courier.NewRouter(ExchangeSessionKey{}))
}

// ExchangeSessionKey 通过Code交换SessionKey
type ExchangeSessionKey struct {
	httpx.MethodGet
	// jsCode
	Code string `in:"query" name:"code"`
}

func (req ExchangeSessionKey) Path() string {
	return "/exchangeSessionKey"
}

func (req ExchangeSessionKey) Output(ctx context.Context) (result interface{}, err error) {
	sessionResp, err := wechat.GetController().Code2Session(req.Code)
	if err != nil {
		return
	}

	userCtrl := user.GetController()
	u, dbErr := userCtrl.GetUserByOpenID(sessionResp.OpenID)
	if u != nil {
		err = userCtrl.UpdateUserInfo(u.UserID, user.UpdateUserInfoParams{
			SessionKey: sessionResp.SessionKey,
		})
		if err != nil {
			return nil, err
		}
		// 刷新token并返回用户信息实体
		return userCtrl.RefreshToken(u.UserID)
	} else {
		if dbErr == errors.UserNotFound {
			// 用户不存在则创建用户
			u, dbErr = userCtrl.CreateUserByWechatSession(user.CreateUserByWechatSessionParams{
				OpenID:     sessionResp.OpenID,
				UnionID:    sessionResp.UnionID,
				SessionKey: sessionResp.SessionKey,
			})
			if dbErr != nil {
				err = dbErr
				return
			}
		} else {
			// 数据库错误
			err = dbErr
			return
		}
	}

	result = u
	return
}
