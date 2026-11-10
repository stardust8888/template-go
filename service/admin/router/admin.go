package router

import (
	"github.com/gin-gonic/gin"

	"template-go/service/admin/controller/admin"
)

func registerAdminRouter(router *gin.RouterGroup) {
	adminRouter := router.Group("/admin")
	adminRouter.POST("/key/generate", admin.Generate)
}
