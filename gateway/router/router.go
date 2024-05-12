package router

import (
	"common/config"
	"gateway/api"
	"github.com/gin-gonic/gin"
)

// RegisterRouter 注册路由
func RegisterRouter() *gin.Engine {
	if config.Conf.Log.Level == "DEBUG" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// 初始化grpc的client gateway是作为grpc的客户端，去调用user grpc服务
	r := gin.Default()
	userHandler := api.NewUserHandler()
	r.POST("/register", userHandler.Register)
	return r
}
