package settlement_flow

import (
	"github.com/eden-framework/sqlx"
	"github.com/eden-w2w/srv-w2w/internal/contants/enums"
	"github.com/eden-w2w/srv-w2w/internal/contants/errors"
	"github.com/eden-w2w/srv-w2w/internal/databases"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/eden-w2w/srv-w2w/internal/modules/id_generator"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var controller *Controller

func GetController() *Controller {
	if controller == nil {
		controller = newController(global.Config.MasterDB, global.Config.SettlementConfig)
	}
	return controller
}

type Controller struct {
	db     sqlx.DBExecutor
	config global.SettlementConfig
	task   *cron.Cron
}

func newController(db sqlx.DBExecutor, config global.SettlementConfig) *Controller {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	t := cron.New(cron.WithParser(parser))
	controller := &Controller{db: db, config: config, task: t}
	_, err := t.AddFunc(config.ToSettlementCronRule(), controller.TaskSettlement)
	if err != nil {
		logrus.Panicf("[settlement_flow.newController] t.AddFunc err: %v, rules: %s", err, config.ToSettlementCronRule())
	}

	return controller
}

func (c Controller) TaskSettlement() {

}

func (c Controller) GetSettlementByUserIDAndName(userID uint64, name string, db sqlx.DBExecutor, forUpdate bool) (model *databases.SettlementFlow, err error) {
	if db == nil {
		db = c.db
	}
	model = &databases.SettlementFlow{
		UserID: userID,
		Name:   name,
	}

	if forUpdate {
		err = model.FetchByUserIDAndNameForUpdate(db)
	} else {
		err = model.FetchByUserIDAndName(db)
	}
	if err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, errors.NotFound
		}
		logrus.Errorf("[GetSettlementByUserIDAndName] err: %v, userID: %d, name: %s", err, userID, name)
		return nil, err
	}

	return
}

func (c Controller) GetPromotionSettlementAmount(flows []databases.PromotionFlow) (totalSales, expectedAmount uint64) {
	for _, flow := range flows {
		totalSales += flow.Amount
	}
	for _, r := range c.config.SettlementRules {
		if totalSales >= r.MinSales && totalSales < r.MaxSales {
			expectedAmount = uint64(float64(totalSales) * r.Proportion)
			return
		}
	}
	return
}

func (c Controller) CreateSettlement(params CreateSettlementParams, db sqlx.DBExecutor) (*databases.SettlementFlow, error) {
	if db == nil {
		db = c.db
	}
	id, _ := id_generator.GetGenerator().GenerateUniqueID()
	model := &databases.SettlementFlow{
		SettlementID: id,
		UserID:       params.UserID,
		Name:         params.Name,
		TotalSales:   params.TotalSales,
		Proportion:   params.Proportion,
		Amount:       params.Amount,
		Status:       enums.SETTLEMENT_STATUS__CREATED,
	}
	err := model.Create(db)
	if err != nil {
		logrus.Errorf("[CreateSettlement] model.Create err: %v, params: %+v", err, params)
		return nil, err
	}
	return model, nil
}
