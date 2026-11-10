package router

import (
	"github.com/gin-gonic/gin"

	"template-go/library/ginx"
	"template-go/service/cron/controller/cron"
	"template-go/utils/constant"
)

func RegisterRoutes(engine *gin.Engine) {
	//engine.Use(ginx.TimeoutMiddleware(30))
	//engine.Use(ginx.LoggerMiddleware())
	engine.Use(ginx.CORSMiddleware())
	engine.Use(ginx.RequestIDMiddleware())
	engine.Use(ginx.RecoveryMiddleware())

	router := engine.Group(constant.BaseCronUrlV1)
	router.GET("/health", ginx.Check)

	cronRouter := router.Group("cron")
	cronRouter.GET("/task", cron.HandleCronTask)
}
