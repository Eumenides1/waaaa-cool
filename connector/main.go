package main

import (
	"common/config"
	"common/logs"
	"common/metrics"
	"connector/app"
	"context"
	"flag"
	"fmt"
	"framework/game"
	"os"
)

var configFile = flag.String("config", "application.yml", "config file")

func main() {
	// 1. 加载配置
	flag.Parse()
	config.InitConfig(*configFile)
	game.InitConfig("../config")
	// 2. 启动监控
	go func() {
		err := metrics.Serve(fmt.Sprintf("0.0.0.0:%d", config.Conf.MetricPort))
		if err != nil {
			panic(err)
		}
	}()
	// 3. 启动grpc服务
	err := app.Run(context.Background(), "connector001")
	if err != nil {
		logs.Error("app run err %v", err)
		os.Exit(-1)
	}
}
