package app

import (
	"common/config"
	"common/logs"
	"context"
	"fmt"
	"gateway/router"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run 启动程序，启动grpc服务，启动Http服务，加载日志 加载数据库
func Run(ctx context.Context) error {
	logs.InitLog(config.Conf.AppName)
	go func() {
		// gin 启动 注册路由
		r := router.RegisterRouter()
		if err := r.Run(fmt.Sprintf(":%d", config.Conf.HttpPort)); err != nil {
			logs.Fatal("gate gin run err :%v", err)
		}
	}()
	// 期望有一个优雅启动和停机
	stop := func() {
		// other
		time.Sleep(3 * time.Second)
		logs.Info("stop app finish")
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
	for {
		select {
		case <-ctx.Done():
			stop()
			return nil
		case s := <-c:
			switch s {
			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
				stop()
				logs.Info("user app quit")
				return nil
			case syscall.SIGHUP:
				stop()
				logs.Info("hang up !! user app quit")
				return nil
			default:
				return nil
			}
		}
	}
}
