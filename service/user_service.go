package service

import (
	"fmt"
	"strings"

	"github.com/Conflux-Chain/go-conflux-util/api"
	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/shopspring/decimal"
	"github.com/v3-Swampy/points-service/model"
	"gorm.io/gorm"
)

type UserService struct {
	store *store.Store
}

func NewUserService(store *store.Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (service *UserService) Get(address string) (*model.User, error) {
	var user model.User
	found, err := service.store.Get(&user, "address = ?", address)
	if err != nil {
		return nil, api.ErrDatabaseCause(err, "Failed to get user by address")
	}

	if !found {
		return nil, api.ErrValidationStr("Failed to find user by address")
	}

	return &user, nil
}

func (service *UserService) DeltaUpsert(address string, tradePoints decimal.Decimal, liquidityPoints decimal.Decimal) error {
	var user model.User
	found, err := service.store.Get(&user, "address = ?", address)
	if err != nil {
		return api.ErrDatabaseCause(err, "Failed to get user by address")
	}

	if !found {
		bean := &model.User{
			Address:         address,
			TradePoints:     tradePoints,
			LiquidityPoints: liquidityPoints,
		}
		return service.store.DB.Create(bean).Error
	}

	newParam := map[string]any{}
	if !tradePoints.Equal(decimal.Zero) {
		newParam["trade_points"] = user.TradePoints.Add(tradePoints)
	}
	if !liquidityPoints.Equal(decimal.Zero) {
		newParam["liquidity_points"] = user.LiquidityPoints.Add(liquidityPoints)
	}

	return service.store.DB.Model(&model.User{}).
		Where("id = ?", user.ID).
		Updates(newParam).Error
}

func (service *UserService) BatchDeltaUpsert(users []*model.User, dbTx ...*gorm.DB) error {
	db := service.store.DB
	if len(dbTx) > 0 {
		db = dbTx[0]
	}

	var placeholders string
	var params []interface{}
	size := len(users)
	for i, u := range users {
		placeholders += "(?,?,?,?,?)"
		if i != size-1 {
			placeholders += ",\n\t\t\t"
		}
		params = append(params, []interface{}{u.Address, u.TradePoints, u.LiquidityPoints, u.CreatedAt, u.UpdatedAt}...)
	}

	sqlString := fmt.Sprintf(`
		insert into 
    		users(address, trade_points, liquidity_points, created_at, updated_at)
		values
			%s
		on duplicate key update
			address = values(address),
			trade_points = trade_points + values(trade_points),
			liquidity_points = liquidity_points + values(liquidity_points),
			updated_at = values(updated_at)                 
	`, placeholders)

	return db.Exec(sqlString, params...).Error
}

func (service *UserService) List(request model.UserPagingRequest) (total int64, users []*model.User, err error) {
	db := service.store.DB.Model(&model.User{})

	if err = db.Count(&total).Error; err != nil {
		return 0, nil, api.ErrDatabaseCause(err, "Failed to get count of users")
	}

	var otherFields string
	if strings.EqualFold(request.SortField, "trade") {
		otherFields = "liquidity_points %s"
	} else {
		otherFields = "trade_points %s"
	}

	var orderBy string
	if request.IsDesc() {
		orderBy = fmt.Sprintf("%s_points DESC, %s", request.SortField, fmt.Sprintf(otherFields, "DESC"))
	} else {
		orderBy = fmt.Sprintf("%s_points ASC, %s", request.SortField, fmt.Sprintf(otherFields, "ASC"))
	}

	if err = db.Order(orderBy).Offset(request.Offset).Limit(request.Limit).Find(&users).Error; err != nil {
		return 0, nil, api.ErrDatabaseCause(err, "Failed to get users")
	}

	return
}
