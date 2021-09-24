package payment

import "github.com/eden-framework/courier"

var Router = courier.NewRouter(PaymentRouter{})

type PaymentRouter struct {
	courier.EmptyOperator
}

func (PaymentRouter) Path() string {
	return "/payment"
}
