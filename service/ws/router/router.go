package router

import (
	"github.com/gin-gonic/gin"

	"template-go/library/ginx"
	"template-go/service/ws/controller/chat"
	"template-go/utils/constant"
)

func RegisterRoutes(engine *gin.Engine) {
	//engine.Use(ginx.TimeoutMiddleware(30))
	//engine.Use(ginx.LoggerMiddleware())
	engine.Use(ginx.CORSMiddleware())
	engine.Use(ginx.RequestIDMiddleware())
	engine.Use(ginx.RecoveryMiddleware())

	websocketRouter := engine.Group(constant.BaseWsUrlV1)
	websocketRouter.GET("/health", ginx.Check)

	websocketRouter.GET("/handshake", chat.Handshake)
	websocketRouter.GET("/connect/count", chat.GetConnectCount)
	websocketRouter.POST("/message/push", chat.MessagePush)
}
