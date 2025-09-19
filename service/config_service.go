package service

import (
	"encoding/json"
	"time"

	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/pkg/errors"
	"github.com/v3-Swampy/points-service/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	CfgKeyLastStatTimePoints = "last.stat.time.points"
	CfxKeyPrefixPoolWeight   = "pool.weight."
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

func (cs *ConfigService) StoreConfig(confName string, confVal interface{}, dbTx ...*gorm.DB) error {
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
		Value: confVal.(string),
	}).Error
}

func (cs *ConfigService) DeleteConfig(confName string) (bool, error) {
	res := cs.store.DB.Delete(&model.Config{}, "name = ?", confName)
	return res.RowsAffected > 0, res.Error
}

// last stat points time

func (cs *ConfigService) GetLastStatPointsTime() (string, error) {
	var cfg model.Config
	err := cs.store.DB.Where("name = ?", CfgKeyLastStatTimePoints).First(&cfg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return cfg.Value, nil
}

func (cs *ConfigService) UpsertLastStatPointsTime(updateTime time.Time, dbTx ...*gorm.DB) error {
	return cs.StoreConfig(CfgKeyLastStatTimePoints, updateTime.Format(time.RFC3339), dbTx...)
}

// pool weight

type PoolWeight struct {
	Trade     *int8 `json:"trade"`
	Liquidity *int8 `json:"liquidity"`
}

func (cs *ConfigService) GetPoolWeight(poolAddress string) (*PoolWeight, error) {
	var cfg model.Config
	err := cs.store.DB.Where("name = ?", CfxKeyPrefixPoolWeight+poolAddress).First(&cfg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	weight := PoolWeight{}
	if err = json.Unmarshal([]byte(cfg.Value), &weight); err != nil {
		return nil, err
	}

	return &weight, nil
}

func (cs *ConfigService) UpsertPoolWeight(poolAddress string, newWeight PoolWeight) error {
	oldWeight, err := cs.GetPoolWeight(poolAddress)
	if err != nil {
		return err
	}

	var weight PoolWeight
	if oldWeight == nil {
		weight = newWeight
	} else {
		if newWeight.Trade != nil {
			oldWeight.Trade = newWeight.Trade
		}
		if newWeight.Liquidity != nil {
			oldWeight.Liquidity = newWeight.Liquidity
		}
		weight = *oldWeight
	}

	val, err := json.Marshal(weight)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal pool weight")
	}

	return cs.StoreConfig(CfxKeyPrefixPoolWeight+poolAddress, string(val))
}
