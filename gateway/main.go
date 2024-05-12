package gateway

import (
	"common/config"
	"common/logs"
	"common/metrics"
	"context"
	"flag"
	"fmt"
	"gateway/app"
	"os"
)

var configFile = flag.String("config", "application.yml", "config file")

func main() {
	// 1. 加载配置
	flag.Parse()
	config.InitConfig(*configFile)
	fmt.Println(config.Conf)
	// 2. 启动监控
	go func() {
		err := metrics.Serve(fmt.Sprintf("0.0.0.0:%d", config.Conf.MetricPort))
		if err != nil {
			panic(err)
		}
	}()
	// 3. 启动grpc服务
	err := app.Run(context.Background())
	if err != nil {
		logs.Error("app run err %v", err)
		os.Exit(-1)
	}
}
