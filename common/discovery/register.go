package discovery

import (
	"common/config"
	"common/logs"
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// Register 将grpc注册到etcd
// 原理 创建一个租约 将grpc服务信息注册到etcd并且绑定租约
// 如果过了租约时间，etcd会删除存储的信息
// 可以实现心跳，完成续租，如果etcd没有则重新注册
type Register struct {
	etcdCli     *clientv3.Client                        //etcd连接
	leaseId     clientv3.LeaseID                        //租约id
	DialTimeout int                                     //超时时间 秒
	ttl         int64                                   //租约时间 秒
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse // 心跳channel
	info        Server                                  //注册的服务信息
	closeCh     chan struct{}
}

func NewRegister() *Register {
	return &Register{
		DialTimeout: 3,
	}
}

func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

func (r *Register) Register(conf config.EtcdConf) error {
	// 注册信息
	info := Server{
		Name:    conf.Register.Name,
		Addr:    conf.Register.Addr,
		Weight:  conf.Register.Weight,
		Version: conf.Register.Version,
		Ttl:     conf.Register.Ttl,
	}
	// 建立etcd的链接
	var err error
	r.etcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:   conf.Addrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		return err
	}
	r.info = info
	if err = r.register(); err != nil {
		return err
	}
	r.closeCh = make(chan struct{})
	go r.watcher()
	return nil
}

func (r *Register) register() error {
	// 创建租约
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(r.DialTimeout))
	defer cancel()
	var err error
	if err = r.createLease(ctx, r.info.Ttl); err != nil {
		return err
	}
	// 心跳检测
	if r.keepAliveCh, err = r.keepAlive(); err != nil {
		return err
	}

	// 绑定租约
	data, _ := json.Marshal(r.info)
	return r.bindLease(ctx, r.info.BuildRegisterKey(), string(data))
}

// createLease 创建租约
func (r *Register) createLease(ctx context.Context, ttl int64) error {
	grant, err := r.etcdCli.Grant(ctx, ttl)
	if err != nil {
		logs.Error("createLease failed err :%v", err)
		return err
	}
	r.leaseId = grant.ID
	return nil
}

func (r *Register) bindLease(ctx context.Context, key, value string) error {
	_, err := r.etcdCli.Put(ctx, key, value, clientv3.WithLease(r.leaseId))
	if err != nil {
		logs.Error("bindLease failed err :%v", err)
		return err
	}
	return nil
}

// keepAlive 心跳检测
func (r *Register) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	alive, err := r.etcdCli.KeepAlive(context.Background(), r.leaseId)
	if err != nil {
		logs.Error("keepAlive failed err :%v", err)
		return alive, err
	}
	return alive, nil
}

// watcher 续约 新注册
func (r *Register) watcher() {
	// 租约到期 需要检查是否自动注册
	ticker := time.NewTicker(time.Duration(r.info.Ttl) * time.Second)
	for {
		select {
		case <-r.closeCh:
			if err := r.unregister(); err != nil {
				logs.Error("close and unRegister failed err :%v", err)
			}
			// 租约撤销
			if _, err := r.etcdCli.Revoke(context.Background(), r.leaseId); err != nil {
				logs.Error("revoke lease failed err :%v", err)
			}
			if r.etcdCli != nil {
				r.etcdCli.Close()
			}
			logs.Info("unregister etcd...")
		case res := <-r.keepAliveCh:
			if res != nil {
				if err := r.register(); err != nil {
					logs.Error("keepAliveCh register failed err :%v", err)
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					logs.Error("ticker register failed err :%v", err)
				}
			}
		}
	}
}

func (r *Register) unregister() error {
	_, err := r.etcdCli.Delete(context.Background(), r.info.BuildRegisterKey())
	return err
}
