package app

import (
	"common/config"
	"common/discovery"
	"common/logs"
	"context"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run 启动程序，启动grpc服务，启动Http服务，加载日志 加载数据库
func Run(ctx context.Context) error {
	logs.InitLog(config.Conf.AppName)
	register := discovery.NewRegister()
	// 启动grpc服务端
	server := grpc.NewServer()
	go func() {
		lis, err := net.Listen("tcp", config.Conf.Grpc.Addr)
		if err != nil {
			logs.Fatal("user grpc server listen err :%v", err)
		}
		err = register.Register(config.Conf.Etcd)
		if err != nil {
			logs.Fatal("user grpc server register etcd err :%v", err)
		}
		// 阻塞操作
		err = server.Serve(lis)
		if err != nil {
			logs.Fatal("user grpc server run failed err :%v", err)
		}
	}()
	// 期望有一个优雅启动和停机
	stop := func() {
		server.Stop()
		register.Stop()
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
