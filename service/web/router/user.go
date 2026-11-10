package router

import (
	"github.com/gin-gonic/gin"

	"template-go/service/web/controller/user"
)

func registerUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("/user")
	userRouter.GET("/delete/apply", user.GetDeleteApply)
	userRouter.POST("/delete/apply", user.ApplyDelete)
	userRouter.POST("/delete/cancel", user.CancelDelete)
}
