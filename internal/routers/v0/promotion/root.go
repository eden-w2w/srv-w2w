package promotion

import "github.com/eden-framework/courier"

var Router = courier.NewRouter(PromotionRouter{})

type PromotionRouter struct {
	courier.EmptyOperator
}

func (PromotionRouter) Path() string {
	return "/promotion"
}
