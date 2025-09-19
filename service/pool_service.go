package service

import (
	"fmt"

	"github.com/Conflux-Chain/go-conflux-util/api"
	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/v3-Swampy/points-service/model"
)

type PoolService struct {
	store *store.Store
}

func NewPoolService(store *store.Store) *PoolService {
	return &PoolService{
		store: store,
	}
}

func (service *PoolService) List(request model.PoolPagingRequest) (total int64, pools []*model.Pool, err error) {
	db := service.store.DB.Model(&model.Pool{})

	if err = db.Count(&total).Error; err != nil {
		return 0, nil, api.ErrDatabaseCause(err, "Failed to get count of pools")
	}

	var sortField string
	if request.SortField == "tvl" {
		sortField = request.SortField
	} else {
		sortField = fmt.Sprintf("%s_points", request.SortField)
	}

	var orderBy string
	if request.IsDesc() {
		orderBy = fmt.Sprintf("%s DESC", sortField)
	} else {
		orderBy = fmt.Sprintf("%s ASC", sortField)
	}

	if err = db.Order(orderBy).Offset(request.Offset).Limit(request.Limit).Find(&pools).Error; err != nil {
		return 0, nil, api.ErrDatabaseCause(err, "Failed to get pools")
	}

	return
}
