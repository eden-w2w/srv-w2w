package payment_flow

import (
	"github.com/eden-framework/sqlx"
	"github.com/eden-framework/sqlx/datatypes"
	"github.com/eden-w2w/srv-w2w/internal/contants/enums"
	"github.com/eden-w2w/srv-w2w/internal/contants/errors"
	"github.com/eden-w2w/srv-w2w/internal/databases"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/eden-w2w/srv-w2w/internal/modules/id_generator"
	"github.com/sirupsen/logrus"
	"time"
)

var controller *Controller

type Controller struct {
	db sqlx.DBExecutor
}

func newController(db sqlx.DBExecutor) *Controller {
	return &Controller{
		db: db,
	}
}

func GetController() *Controller {
	if controller == nil {
		controller = newController(global.Config.MasterDB)
	}
	return controller
}

func (c Controller) CreatePaymentFlow(params CreatePaymentFlowParams, db sqlx.DBExecutor) (*databases.PaymentFlow, error) {
	if db == nil {
		db = c.db
	}
	id, _ := id_generator.GetGenerator().GenerateUniqueID()
	model := &databases.PaymentFlow{
		FlowID:        id,
		UserID:        params.UserID,
		OrderID:       params.OrderID,
		Amount:        params.Amount,
		PaymentMethod: params.PaymentMethod,
		Status:        enums.PAYMENT_STATUS__CREATED,
		ExpiredAt:     datatypes.MySQLTimestamp(time.Now().Add(global.Config.PaymentFlowExpireIn)),
	}
	err := model.Create(c.db)
	if err != nil {
		logrus.Errorf("[CreatePaymentFlow] model.Create err: %v, params: %+v", err, params)
		return nil, errors.InternalError
	}
	return model, nil
}

func (c Controller) GetFlowByOrderAndUserID(orderID, userID uint64, db sqlx.DBExecutor) (*databases.PaymentFlow, error) {
	if db == nil {
		db = c.db
	}

	model := &databases.PaymentFlow{}
	models, err := model.BatchFetchByOrderAndUserID(db, orderID, userID, enums.PAYMENT_STATUS__SUCCESS)
	if err != nil {
		logrus.Errorf("[GetFlowByOrderAndUserID] model.BatchFetchByOrderAndUserID err: %v, orderID: %d, userID: %d", err, orderID, userID)
		return nil, errors.InternalError
	}

	if len(models) == 0 {
		logrus.Errorf("[GetFlowByOrderAndUserID] len(models) == 0, orderID: %d, userID: %d", orderID, userID)
		return nil, errors.PaymentFlowNotFound
	}

	return &models[0], nil
}
