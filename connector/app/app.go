package app

import (
	"common/config"
	"common/logs"
	"context"
	"framework/connector"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run 启动程序，启动grpc服务，启动Http服务，加载日志 加载数据库
func Run(ctx context.Context, serverId string) error {
	logs.InitLog(config.Conf.AppName)
	exit := func() {}
	go func() {
		c := connector.Default()
		exit = c.Close
		c.Run(serverId)
	}()
	// 期望有一个优雅启动和停机
	stop := func() {
		// other
		exit()
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
				logs.Info("connector app quit")
				return nil
			case syscall.SIGHUP:
				stop()
				logs.Info("hang up !! connector app quit")
				return nil
			default:
				return nil
			}
		}
	}
}
