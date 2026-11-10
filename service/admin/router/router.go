package router

import (
	"github.com/gin-gonic/gin"

	"template-go/library/ginx"
	"template-go/service/admin/controller/admin"
	"template-go/utils/constant"
)

func RegisterRoutes(engine *gin.Engine) {
	//engine.Use(ginx.TimeoutMiddleware(30))
	//engine.Use(ginx.LoggerMiddleware())
	engine.Use(ginx.CORSMiddleware())
	engine.Use(ginx.RequestIDMiddleware())
	engine.Use(ginx.RecoveryMiddleware())

	router := engine.Group(constant.BaseAdminUrlV1)
	router.GET("/health", ginx.Check)
	router.POST("/administrator/login", admin.Login)
	router.Use(ginx.AdminAuthMiddleware())
	registerAdminRouter(router)
	registerUserRouter(router)
	registerFeedbackRouter(router)
}
