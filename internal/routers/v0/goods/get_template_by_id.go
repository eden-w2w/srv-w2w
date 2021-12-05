package goods

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-w2w/lib-modules/modules/freight_template"
)

func init() {
	Router.Register(courier.NewRouter(GetTemplateById{}))
}

// GetTemplateById 通过ID获取运费模板
type GetTemplateById struct {
	httpx.MethodGet
	TemplateID uint64 `in:"path" name:"templateID,string"`
}

func (req GetTemplateById) Path() string {
	return "/:goodsID/freight/:templateID"
}

func (req GetTemplateById) Output(ctx context.Context) (result interface{}, err error) {
	return freight_template.GetController().GetTemplateByID(req.TemplateID, nil, false)
}
