package settings

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-w2w/lib-modules/modules/settings"
)

func init() {
	Router.Register(courier.NewRouter(GetSettings{}))
}

// GetSettings 获取系统配置
type GetSettings struct {
	httpx.MethodGet
}

func (req GetSettings) Path() string {
	return ""
}

func (req GetSettings) Output(ctx context.Context) (result interface{}, err error) {
	setting := settings.GetController().GetSetting()
	return setting, nil
}
