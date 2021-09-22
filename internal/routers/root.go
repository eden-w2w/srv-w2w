package routers

import (
	"github.com/eden-framework/courier"
	v0 "github.com/eden-w2w/srv-w2w/internal/routers/v0"
)

var Router = courier.NewRouter(RootRouter{})

func init() {
	Router.Register(v0.Router)
}

type RootRouter struct {
	courier.EmptyOperator
}

func (RootRouter) Path() string {
	return "/w2w"
}
