package user

import (
	"crypto"
	"encoding/hex"
	"github.com/eden-framework/sqlx"
	"github.com/eden-framework/sqlx/datatypes"
	"github.com/eden-w2w/srv-w2w/internal"
	"github.com/eden-w2w/srv-w2w/internal/contants/errors"
	"github.com/eden-w2w/srv-w2w/internal/databases"
	"github.com/eden-w2w/srv-w2w/internal/global"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
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

func (c Controller) CreateUserByWechatSession(params CreateUserByWechatSessionParams) (*databases.User, error) {
	id, _ := internal.Generator.GenerateUniqueID()
	model := &databases.User{
		UserID:      id,
		Token:       c.generateToken(id),
		OpenID:      params.OpenID,
		UnionID:     params.UnionID,
		SessionKey:  params.SessionKey,
		OperateTime: datatypes.OperateTime{},
	}
	err := model.Create(c.db)
	if err != nil {
		logrus.Errorf("[CreateUserByWechatSession] model.Create(c.db) err: %v, params: %+v", err, params)
		return nil, errors.InternalError
	}

	return model, nil
}

func (c Controller) RefreshToken(userID uint64) (*databases.User, error) {
	token := c.generateToken(userID)
	model := &databases.User{
		UserID: userID,
		Token:  token,
	}
	err := model.UpdateByUserIDWithStruct(c.db)
	if err != nil {
		logrus.Errorf("[RefreshToken] model.UpdateByUserIDWithStruct(c.db) err: %v, userID: %d", err, userID)
		return nil, errors.InternalError
	}

	err = model.FetchByUserID(c.db)
	if err != nil {
		logrus.Errorf("[RefreshToken] model.UpdateByUserIDWithStruct(c.db) err: %v, userID: %d", err, userID)
		return nil, errors.InternalError
	}

	return model, nil
}

func (c Controller) UpdateUserInfo(userID uint64, params UpdateUserInfoParams) error {
	model := databases.User{
		UserID: userID,
	}
	err := model.FetchByUserID(c.db)
	if err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return errors.UserNotFound
		}
		logrus.Errorf("[UpdateUserInfo] model.FetchByUserID(c.db) err: %v, userID: %d, params: %+v", err, userID, params)
		return err
	}

	if params.Diff(&model) {
		err = model.UpdateByUserIDWithStruct(c.db)
		if err != nil {
			logrus.Errorf("[UpdateUserInfo] model.UpdateByUserIDWithStruct(c.db) err: %v, userID: %d, params: %+v", err, userID, params)
			return errors.InternalError
		}
	}

	return nil
}

func (c Controller) generateToken(userID uint64) string {
	id := strconv.FormatUint(userID, 10)
	t := strconv.FormatInt(time.Now().UnixNano(), 10)
	sha256 := crypto.SHA256.New()
	sha256.Write([]byte(id + t))
	hash := sha256.Sum(nil)
	return hex.EncodeToString(hash)
}

func (c Controller) GetUserByOpenID(openID string) (*databases.User, error) {
	model := &databases.User{
		OpenID: openID,
	}
	err := model.FetchByOpenID(c.db)
	if err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, errors.UserNotFound
		}
		logrus.Errorf("[GetUserByOpenID] model.FetchByOpenID err: %v, openID: %s", err, openID)
		return nil, errors.InternalError
	}
	return model, nil
}

func (c Controller) GetUserByToken(token string) (*databases.User, error) {
	model := &databases.User{
		Token: token,
	}
	err := model.FetchByToken(c.db)
	if err != nil {
		if sqlx.DBErr(err).IsNotFound() {
			return nil, errors.UserNotFound
		}
		logrus.Errorf("[GetUserByToken] model.FetchByToken err: %v, token: %s", err, token)
		return nil, errors.InternalError
	}
	return model, nil
}
