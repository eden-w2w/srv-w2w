package goods

import (
	"context"
	"github.com/eden-framework/courier"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-w2w/lib-modules/constants/enums"
	"github.com/eden-w2w/lib-modules/modules/booking_flow"
)

func init() {
	Router.Register(courier.NewRouter(GetGoodsBooking{}))
}

// GetGoodsBooking 获取商品预售状态
type GetGoodsBooking struct {
	httpx.MethodGet
	// 商品ID
	GoodsID uint64 `in:"path" name:"goodsID"`
}

func (req GetGoodsBooking) Path() string {
	return "/:goodsID/booking"
}

type GetGoodsBookingResponse struct {
	Booking bool `json:"booking"`
}

func (req GetGoodsBooking) Output(ctx context.Context) (result interface{}, err error) {
	flows, err := booking_flow.GetController().GetBookingFlowByGoodsIDAndStatus(
		req.GoodsID,
		enums.BOOKING_STATUS__PROCESS,
	)
	if err != nil {
		return
	}
	if len(flows) == 0 {
		return &GetGoodsBookingResponse{Booking: false}, nil
	}
	return &GetGoodsBookingResponse{Booking: true}, nil
}
