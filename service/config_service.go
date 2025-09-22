package service

import (
	"strconv"

	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/pkg/errors"
	"github.com/v3-Swampy/points-service/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	CfgKeyLastStatTimePoints = "last.stat.time.points"
)

type ConfigService struct {
	store *store.Store
}

func NewConfigService(store *store.Store) *ConfigService {
	return &ConfigService{
		store: store,
	}
}

func (cs *ConfigService) LoadConfig(confNames ...string) (map[string]interface{}, error) {
	var confs []model.Config

	if err := cs.store.DB.Where("name IN ?", confNames).Find(&confs).Error; err != nil {
		return nil, err
	}

	res := make(map[string]interface{}, len(confs))
	for _, c := range confs {
		res[c.Name] = c.Value
	}

	return res, nil
}

func (cs *ConfigService) StoreConfig(confName string, confVal string, dbTx ...*gorm.DB) error {
	db := cs.store.DB
	if len(dbTx) > 0 {
		db = dbTx[0]
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"value":      confVal,
			"updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}),
	}).Create(&model.Config{
		Name:  confName,
		Value: confVal,
	}).Error
}

func (cs *ConfigService) DeleteConfig(confName string) (bool, error) {
	res := cs.store.DB.Delete(&model.Config{}, "name = ?", confName)
	return res.RowsAffected > 0, res.Error
}

// last stat points time

func (cs *ConfigService) GetLastStatPointsTime() (int64, error) {
	var cfg model.Config
	err := cs.store.DB.Where("name = ?", CfgKeyLastStatTimePoints).First(&cfg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	timestamp, err := strconv.ParseInt(cfg.Value, 10, 64)
	if err != nil {
		return 0, err
	}

	return timestamp, nil
}

func (cs *ConfigService) UpsertLastStatPointsTime(timestamp int64, dbTx ...*gorm.DB) error {
	return cs.StoreConfig(CfgKeyLastStatTimePoints, strconv.FormatInt(timestamp, 10), dbTx...)
}
