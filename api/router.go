package api

import (
	"github.com/Conflux-Chain/go-conflux-util/api"
	"github.com/Conflux-Chain/go-conflux-util/api/middleware"
	"github.com/Conflux-Chain/go-conflux-util/viper"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/v3-Swampy/points-service/docs"
	"github.com/v3-Swampy/points-service/service"
)

func MustServeFromViper(services service.Services) {
	var config api.Config
	viper.MustUnmarshalKey("api", &config)

	api.MustServe(config, func(router *gin.Engine) {
		Routes(router, services)
	})
}

//	@title			Points Service API
//	@version		1.0
//	@description	Use any http client to fetch data from the Points Service

func Routes(router *gin.Engine, services service.Services) {
	docs.SwaggerInfo.BasePath = "/api"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	controller := NewController(services)

	router.GET("/api/users", middleware.Wrap(controller.listUsers))
	router.GET("/api/pools", middleware.Wrap(controller.listPools))

	logrus.Info("Service started")
}
