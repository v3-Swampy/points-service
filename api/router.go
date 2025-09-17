package api

import (
	"github.com/Conflux-Chain/go-conflux-util/api/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.GET("/api/users", middleware.Wrap(listUsers))
	router.GET("/api/pools", middleware.Wrap(listPools))
}
