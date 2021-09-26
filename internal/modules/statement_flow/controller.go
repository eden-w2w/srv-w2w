package statement_flow

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
		controller = newController(global.Config.MasterDB, global.Config.StatementConfig)
	}
	return controller
}

type Controller struct {
	db     sqlx.DBExecutor
	config global.StatementConfig
	task   *cron.Cron
}

func newController(db sqlx.DBExecutor, config global.StatementConfig) *Controller {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	t := cron.New(cron.WithParser(parser))
	controller := &Controller{db: db, config: config, task: t}
	_, err := t.AddFunc(config.ToStatementCronRule(), controller.TaskStatements)
	if err != nil {
		logrus.Panicf("[statement_flow.newController] t.AddFunc err: %v, rules: %s", err, config.ToStatementCronRule())
	}

	return controller
}

func (c Controller) TaskStatements() {

}

func (c Controller) GetStatementByUserIDAndName(userID uint64, name string, db sqlx.DBExecutor, forUpdate bool) (model *databases.StatementsFlow, err error) {
	if db == nil {
		db = c.db
	}
	model = &databases.StatementsFlow{
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
		logrus.Errorf("[GetStatementByUserIDAndName] err: %v, userID: %d, name: %s", err, userID, name)
		return nil, err
	}

	return
}

func (c Controller) CreateStatement(params CreateStatementParams, db sqlx.DBExecutor) (*databases.StatementsFlow, error) {
	if db == nil {
		db = c.db
	}
	id, _ := id_generator.GetGenerator().GenerateUniqueID()
	model := &databases.StatementsFlow{
		StatementsID: id,
		UserID:       params.UserID,
		Name:         params.Name,
		TotalSales:   params.TotalSales,
		Proportion:   params.Proportion,
		Amount:       params.Amount,
		Status:       enums.STATEMENT_STATUS__CREATED,
	}
	err := model.Create(db)
	if err != nil {
		logrus.Errorf("[CreateStatement] model.Create err: %v, params: %+v", err, params)
		return nil, err
	}
	return model, nil
}
