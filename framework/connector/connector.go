package connector

import (
	"common/logs"
	"fmt"
	"framework/game"
	"framework/net"
)

type Connector struct {
	isRunning bool
	wsManager *net.Manager
}

func Default() *Connector {
	return &Connector{}
}

func (c *Connector) Run(serverId string) {
	if !c.isRunning {
		// 启动websocket和nats
		c.Serve(serverId)
	}
}

func (c *Connector) Close() {
	if c.isRunning {
		// 关闭websocket和nats
	}
}

func (c *Connector) Serve(serverId string) {
	connectorConfig := game.Conf.GetConnector(serverId)
	if connectorConfig == nil {
		logs.Fatal("no connector config found")
	}
	addr := fmt.Sprintf("%s:%d", connectorConfig.Host, connectorConfig.ClientPort)
	c.wsManager.Run(addr)
}
