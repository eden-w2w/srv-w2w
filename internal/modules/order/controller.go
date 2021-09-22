package order

import (
	"github.com/eden-framework/sqlx"
	"github.com/eden-framework/sqlx/builder"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/sirupsen/logrus"

	"github.com/eden-w2w/srv-w2w/internal"
	"github.com/eden-w2w/srv-w2w/internal/contants/enums"
	"github.com/eden-w2w/srv-w2w/internal/contants/errors"
	"github.com/eden-w2w/srv-w2w/internal/databases"
)

var controller = NewController(global.Config.MasterDB)

func GetController() *Controller {
	return controller
}

type Controller struct {
	db sqlx.DBExecutor
}

func NewController(db sqlx.DBExecutor) *Controller {
	return &Controller{db: db}
}

func (c Controller) CreateOrder(p CreateOrderParams, locker InventoryLock) (*databases.Order, error) {
	// 获取订单总额与库中物料进行比对
	var totalPrice uint64 = 0
	var goodsList = make([]CreateOrderGoodsModelParams, 0)
	for _, g := range p.Goods {
		goods := databases.Goods{GoodsID: g.GoodsID}
		err := goods.FetchByGoodsID(c.db)
		if err != nil {
			logrus.Errorf("[CreateOrder] goods.FetchByGoodsID(c.db) err: %v, goodsID: %d", err, g.GoodsID)
			return nil, errors.NotFound.StatusError().WithDesc("商品无法找到")
		}
		totalPrice += goods.Price * uint64(g.Amount)
		goodsList = append(goodsList, CreateOrderGoodsModelParams{
			Goods:  goods,
			Amount: g.Amount,
		})
	}
	if totalPrice != p.TotalPrice {
		logrus.Errorf("[CreateOrder] totalPrice != p.TotalPrice totalPrice: %d, p.TotalPrice: %d", totalPrice, p.TotalPrice)
		return nil, errors.BadRequest.StatusError().WithDesc("订单总额与商品总价不一致")
	}
	if len(goodsList) == 0 {
		logrus.Errorf("[CreateOrder] len(goodsList) == 0")
		return nil, errors.BadRequest.StatusError().WithDesc("商品列表为空")
	}

	// 创建订单
	id, _ := internal.Generator.GenerateUniqueID()
	order := &databases.Order{
		OrderID:       id,
		UserID:        p.UserID,
		RefererID:     p.RefererID,
		TotalPrice:    p.TotalPrice,
		PaymentMethod: p.PaymentMethod,
		Remark:        p.Remark,
		Recipients:    p.Recipients,
		ShippingAddr:  p.ShippingAddr,
		Mobile:        p.Mobile,
		Status:        enums.ORDER_STATUS__CREATED,
	}

	tx := sqlx.NewTasks(c.db)
	tx = tx.With(func(db sqlx.DBExecutor) error {
		return order.Create(db)
	})

	// 锁定库存
	tx = tx.With(func(db sqlx.DBExecutor) error {
		for _, item := range goodsList {
			err := locker(db, item.GoodsID, item.Amount)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// 创建订单物料
	tx = tx.With(func(db sqlx.DBExecutor) error {
		for _, item := range goodsList {
			orderGoods := &databases.OrderGoods{
				OrderID:        id,
				GoodsID:        item.GoodsID,
				Name:           item.Name,
				Comment:        item.Comment,
				DispatchAddr:   item.DispatchAddr,
				Sales:          item.Sales,
				MainPicture:    item.MainPicture,
				Pictures:       item.Pictures,
				Specifications: item.Specifications,
				Activities:     item.Activities,
				LogisticPolicy: item.LogisticPolicy,
				Price:          item.Price,
				Inventory:      item.Inventory,
				Detail:         item.Detail,
				Amount:         item.Amount,
			}
			err := orderGoods.Create(db)
			if err != nil {
				return err
			}
		}
		return nil
	})
	err := tx.Do()
	if err != nil {
		logrus.Errorf("[CreateOrder] err: %v, params: %+v", err, p)
		return nil, errors.InternalError
	}

	return order, nil
}

func (c Controller) GetOrder(orderID, userID uint64) (*databases.Order, error) {
	order := &databases.Order{OrderID: orderID}
	err := order.FetchByOrderID(c.db)
	if err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, errors.OrderNotFound
		}
		logrus.Errorf("[GetOrder] err: %v, orderID: %d", err, orderID)
		return nil, errors.InternalError
	}
	if order.UserID != userID {
		logrus.Errorf("[GetOrder] order.UserID != userID, order.UserID: %d, userID: %d", order.UserID, userID)
		return nil, errors.Forbidden
	}
	return order, nil
}

func (c Controller) GetOrders(p GetOrdersParams) ([]databases.Order, error) {
	order := databases.Order{}
	orders, err := order.List(c.db, p.Conditions(), p.Additions()...)
	if err != nil {
		logrus.Errorf("[GetOrders] order.List err: %v, params: %+v", err, p)
		return nil, errors.InternalError
	}
	return orders, nil
}

func (c Controller) GetOrderGoods(orderID uint64) ([]databases.OrderGoods, error) {
	og := databases.OrderGoods{}
	goods, err := og.BatchFetchByOrderIDList(c.db, []uint64{orderID})
	if err != nil {
		logrus.Errorf("[GetOrderGoods] og.BatchFetchByOrderIDList err: %v, orderID: %d", err, orderID)
		return nil, errors.InternalError
	}
	return goods, nil
}

func (c Controller) UpdateOrderStatusWithDB(db sqlx.DBExecutor, orderID uint64, status enums.OrderStatus) error {
	order := &databases.Order{OrderID: orderID}
	err := order.FetchByOrderID(db)
	if err != nil {
		logrus.Errorf("[UpdateOrderStatus] order.FetchByOrderID err: %v, orderID: %d, status: %s", err, orderID, status.String())
		return errors.InternalError
	}

	if order.Status == status {
		return nil
	}

	// TODO 状态流转检查
	f := builder.FieldValues{
		"Status": status,
	}
	err = order.UpdateByIDWithMap(c.db, f)
	if err != nil {
		logrus.Errorf("[UpdateOrderStatus] order.UpdateByIDWithStruct err: %v, orderID: %d, status: %s", err, orderID, status.String())
		return errors.InternalError
	}
	return nil
}

func (c Controller) UpdateOrderStatus(orderID uint64, status enums.OrderStatus) error {
	return c.UpdateOrderStatusWithDB(c.db, orderID, status)
}

func (c Controller) CancelOrder(orderID, userID uint64, unlocker InventoryUnlock) error {
	order, err := c.GetOrder(orderID, userID)
	if err != nil {
		return err
	}

	if order.Status == enums.ORDER_STATUS__CANCEL {
		logrus.Errorf("[CancelOrder] order.Status == enums.ORDER_STATUS__CANCEL, order.UserID=%d, userID=%d", order.UserID, userID)
		return errors.OrderCanceled
	}

	tx := sqlx.NewTasks(c.db)
	tx = tx.With(func(db sqlx.DBExecutor) error {
		return c.UpdateOrderStatusWithDB(db, orderID, enums.ORDER_STATUS__CANCEL)
	})

	tx = tx.With(func(db sqlx.DBExecutor) error {
		goods, err := c.GetOrderGoods(orderID)
		if err != nil {
			return err
		}

		for _, g := range goods {
			err = unlocker(db, g.GoodsID, g.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	err = tx.Do()
	if err != nil {
		logrus.Errorf("[CancelOrder] tx.Do() err: %v", err)
	}
	return err
}
