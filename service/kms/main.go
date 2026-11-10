package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"template-go/conf"
	"time"

	"github.com/gin-gonic/gin"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"

	"template-go/library/security"
	"template-go/logs"
	"template-go/service/kms/router"
)

func main() {
	var (
		ctx = context.WithValue(context.Background(), logs.RequestIDKey, "start...")
		err error
	)

	if err = conf.InitKmsConfig(); err != nil {
		panic("初始化配置文件失败:" + err.Error())
	}

	// 加载主密钥：环境变量 KMS_MASTER_KEY 优先，缺失时自动生成并落盘。
	// 加载失败直接拒绝启动——带着不可用的密钥继续运行只会让加密结果失去意义。
	if err = security.Init(); err != nil {
		panic("加载主密钥失败:" + err.Error())
	}

	logs.InitLogger("crypto.log")
	defer logs.Sync()

	engine := gin.New()
	router.RegisterRoutes(engine)

	server := &http.Server{
		Addr:    "127.0.0.1:12006", // 仅监听本机回环地址，外部无法访问
		Handler: engine,
	}
	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logs.Fatal(ctx, "HTTP服务器启动失败", zap.Error(err))
		}
	}()
	logs.Info(ctx, "服务器已启动...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logs.Info(ctx, "正在优雅关闭服务器...")

	time.Sleep(3 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		logs.Error(ctx, "服务器优雅关闭失败", zap.Error(err))
	} else {
		logs.Info(ctx, "服务器已安全关闭")
	}

	logs.Info(ctx, "服务已退出")
}
