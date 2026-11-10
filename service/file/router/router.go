package router

import (
	"github.com/gin-gonic/gin"

	"template-go/library/ginx"
	"template-go/service/file/controller/file"
	"template-go/utils/constant"
)

func RegisterRoutes(engine *gin.Engine) {
	//engine.Use(ginx.TimeoutMiddleware(30))
	//engine.Use(ginx.LoggerMiddleware())
	engine.Use(ginx.CORSMiddleware())
	engine.Use(ginx.RequestIDMiddleware())
	engine.Use(ginx.RecoveryMiddleware())

	router := engine.Group(constant.BaseFileUrlV1)
	router.GET("/health", ginx.Check)

	router.Use(ginx.AuthMiddleware())
	fileRouter := router.Group("/ff")
	fileRouter.GET("/view", file.View)
	fileRouter.POST("/view", file.View)
	fileRouter.POST("/upload", file.Upload)
	fileRouter.POST("/delete", file.Delete)
}
