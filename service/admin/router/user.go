package router

import (
	"github.com/gin-gonic/gin"

	"template-go/service/admin/controller/user"
)

func registerUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("/user")
	userRouter.GET("/list", user.GetList)
}
