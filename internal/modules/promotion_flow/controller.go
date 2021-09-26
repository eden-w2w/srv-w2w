package promotion_flow

import (
	"github.com/eden-framework/sqlx"
	"github.com/eden-w2w/srv-w2w/internal/contants/errors"
	"github.com/eden-w2w/srv-w2w/internal/databases"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/eden-w2w/srv-w2w/internal/modules/id_generator"
	"github.com/sirupsen/logrus"
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

func (c Controller) CreatePromotionFlow(params CreatePromotionFlowParams, db sqlx.DBExecutor) (*databases.PromotionFlow, error) {
	if db == nil {
		db = c.db
	}
	id, _ := id_generator.GetGenerator().GenerateUniqueID()
	model := &databases.PromotionFlow{
		FlowID:          id,
		UserID:          params.UserID,
		UserNickName:    params.UserNickName,
		RefererID:       params.RefererID,
		RefererNickName: params.RefererNickName,
		Amount:          params.Amount,
		PaymentFlowID:   params.PaymentFlowID,
	}
	err := model.Create(db)
	if err != nil {
		logrus.Errorf("[CreatePromotionFlow] model.Create err: %v, params: %+v", err, params)
		return nil, errors.InternalError
	}
	return model, nil
}

func (c Controller) GetPromotionFlows(params GetPromotionFlowParams) (list []databases.PromotionFlow, err error) {
	model := &databases.PromotionFlow{}
	list, err = model.List(c.db, params.Conditions(), params.Additions()...)
	if err != nil {
		logrus.Errorf("[GetPromotionFlows] model.List err: %v, params: %+v", err, params)
	}
	return
}
