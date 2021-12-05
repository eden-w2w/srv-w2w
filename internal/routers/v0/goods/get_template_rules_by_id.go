package goods

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-w2w/lib-modules/modules/freight_template"
)

func init() {
	Router.Register(courier.NewRouter(GetTemplateRulesById{}))
}

// GetTemplateRulesById 根据ID获取模板规则
type GetTemplateRulesById struct {
	httpx.MethodGet
	TemplateID uint64 `in:"path" name:"templateID,string"`
}

func (req GetTemplateRulesById) Path() string {
	return "/:goodsID/freight/:templateID/rules"
}

func (req GetTemplateRulesById) Output(ctx context.Context) (result interface{}, err error) {
	return freight_template.GetController().GetTemplateRules(req.TemplateID, freight_template.GetTemplateRuleParams{})
}
