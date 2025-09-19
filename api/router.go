package api

import (
	"github.com/Conflux-Chain/go-conflux-util/api"
	"github.com/Conflux-Chain/go-conflux-util/api/middleware"
	"github.com/Conflux-Chain/go-conflux-util/store"
	"github.com/Conflux-Chain/go-conflux-util/viper"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func MustServeFromViper(store *store.Store) {
	var config api.Config
	viper.MustUnmarshalKey("api", &config)

	api.MustServe(config, func(router *gin.Engine) {
		Routes(router, store)
	})
}

func Routes(router *gin.Engine, store *store.Store) {
	controller := NewController(store)

	router.GET("/api/users", middleware.Wrap(controller.listUsers))
	router.GET("/api/pools", middleware.Wrap(controller.listPools))

	logrus.Info("Service started")
}
