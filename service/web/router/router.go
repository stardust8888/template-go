package router

import (
	"github.com/gin-gonic/gin"

	"template-go/library/ginx"
	"template-go/service/web/controller/user"
	"template-go/utils/constant"
)

func RegisterRoutes(engine *gin.Engine) {
	//engine.Use(ginx.TimeoutMiddleware(30))
	//engine.Use(ginx.LoggerMiddleware())
	engine.Use(ginx.CORSMiddleware())
	engine.Use(ginx.RequestIDMiddleware())
	engine.Use(ginx.RecoveryMiddleware())

	router := engine.Group(constant.BaseWebUrlV1)
	router.GET("/health", ginx.Check)
	router.GET("/user/phone", user.GetUserPhone)
	router.POST("/fast/login", user.FastLogin)

	router.Use(ginx.AuthMiddleware())
	router.POST("/loginOut", user.LoginOut)
	registerUserRouter(router)
	registerFeedbackRouter(router)
}
