package service

import (
	"fmt"

	"github.com/Conflux-Chain/go-conflux-util/api"
	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/v3-Swampy/points-service/model"
)

type UserService struct {
	store *store.Store
}

func NewUserService(store *store.Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (service *UserService) List(request model.UserPagingRequest) (users []*model.User, err error) {
	db := service.store.DB.Model(&model.User{})

	var orderBy string
	if request.IsDesc() {
		orderBy = fmt.Sprintf("%s_points DESC", request.SortField)
	} else {
		orderBy = fmt.Sprintf("%s_points ASC", request.SortField)
	}

	if err = db.Order(orderBy).Offset(request.Offset).Limit(request.Limit).Find(&users).Error; err != nil {
		return nil, api.ErrDatabaseCause(err, "Failed to get users")
	}

	return
}
