package statement

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

func GetController() *Controller {
	if controller == nil {
		controller = newController(global.Config.MasterDB)
	}
	return controller
}

type Controller struct {
	db sqlx.DBExecutor
}

func newController(db sqlx.DBExecutor) *Controller {
	return &Controller{db: db}
}

func (c Controller) GetStatementByInterval(userID uint64, startTime, endTime datatypes.MySQLTimestamp, db sqlx.DBExecutor, forUpdate bool) (model *databases.Statements, err error) {
	if db == nil {
		db = c.db
	}
	model = &databases.Statements{
		UserID:    userID,
		StartTime: startTime,
		EndTime:   endTime,
	}

	if forUpdate {
		err = model.FetchByUserIDAndStartTimeAndEndTimeForUpdate(db)
	} else {
		err = model.FetchByUserIDAndStartTimeAndEndTime(db)
	}
	if err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, errors.NotFound
		}
		logrus.Errorf("[GetStatementByInterval] err: %v, userID: %d, startTime: %s, endTime: %s", err, userID, startTime, endTime)
		return nil, err
	}

	return
}

func (c Controller) CreateStatement(params CreateStatementParams, db sqlx.DBExecutor) (*databases.Statements, error) {
	if db == nil {
		db = c.db
	}
	id, _ := id_generator.GetGenerator().GenerateUniqueID()
	model := &databases.Statements{
		StatementsID: id,
		UserID:       params.UserID,
		StartTime:    params.StartTime,
		EndTime:      params.EndTime,
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

func (c Controller) CreateStatementByNow(userID, amount uint64, db sqlx.DBExecutor) (*databases.Statements, error) {
	if db == nil {
		db = c.db
	}

	var startTime, endTime datatypes.MySQLTimestamp
	if global.Config.StatementType == enums.STATEMENT_TYPE__WEEK {
		startTime, endTime = getWeekTimeIntervalOfThisStatement(global.Config.StatementDate)
	} else if global.Config.StatementType == enums.STATEMENT_TYPE__MONTH {
		startTime, endTime = getMonthTimeIntervalOfThisStatement(global.Config.StatementDate)
	}

	statement, err := c.GetStatementByInterval(userID, startTime, endTime, db, true)
	if err != nil {
		if err == errors.NotFound {
			statement, err = c.CreateStatement(CreateStatementParams{
				UserID:    userID,
				StartTime: startTime,
				EndTime:   endTime,
				Amount:    amount,
			}, db)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		statement.Amount += amount
		err = statement.UpdateByStatementsIDWithStruct(db)
		if err != nil {
			logrus.Errorf("[CreateStatementByNow] statement.UpdateByStatementsIDWithStruct err: %v, statement: %+v", err, statement)
			return nil, err
		}
	}
	return statement, nil
}

func getWeekTimeIntervalOfThisStatement(dayOfWeek uint8) (startTime datatypes.MySQLTimestamp, endTime datatypes.MySQLTimestamp) {
	now := time.Now()
	offset := int(time.Weekday(dayOfWeek) - now.Weekday())
	if offset < 0 {
		offset += 7
	}
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	start := end.AddDate(0, 0, 7)

	startTime = datatypes.MySQLTimestamp(start)
	endTime = datatypes.MySQLTimestamp(end)
	return
}

func getMonthTimeIntervalOfThisStatement(dayOfMonth uint8) (startTime datatypes.MySQLTimestamp, endTime datatypes.MySQLTimestamp) {
	now := time.Now()
	month := now.Month()
	offset := int(dayOfMonth) - now.Day()
	if offset < 0 {
		month++
	}
	end := time.Date(now.Year(), month, int(dayOfMonth), 0, 0, 0, 0, time.Local)
	start := time.Date(now.Year(), month-1, int(dayOfMonth), 0, 0, 0, 0, time.Local)

	startTime = datatypes.MySQLTimestamp(start)
	endTime = datatypes.MySQLTimestamp(end)
	return
}
