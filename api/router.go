package api

import (
	"github.com/Conflux-Chain/go-conflux-util/api"
	"github.com/Conflux-Chain/go-conflux-util/api/middleware"
	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/Conflux-Chain/go-conflux-util/viper"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func MustServeFromViper(db *gorm.DB) {
	var config api.Config
	viper.MustUnmarshalKey("api", &config)

	api.MustServe(config, func(router *gin.Engine) {
		Routes(router, db)
	})
}

func Routes(router *gin.Engine, db *gorm.DB) {
	controller := NewController(store.NewStore(db))

	router.GET("/api/users", middleware.Wrap(controller.listUsers))
	router.GET("/api/pools", middleware.Wrap(controller.listPools))

	logrus.Info("Service started")
}
