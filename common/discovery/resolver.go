package discovery

import (
	"common/config"
	"common/logs"
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"time"
)

type Resolver struct {
	conf        config.EtcdConf
	etcdCli     *clientv3.Client
	DialTimeout int
	closeCh     chan struct{}
	key         string
	cc          resolver.ClientConn
	srvAddrList []resolver.Address
	watchCh     clientv3.WatchChan
}

// Build 当grpc.dial的时候，会同步调用此方法
func (r Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// 获取调用key 链接etcd 获取value
	r.cc = cc
	// 1. 链接etcd
	var err error
	r.etcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.conf.Addrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		logs.Fatal("grpc client connect etcd err : %v", err)
	}
	r.closeCh = make(chan struct{})
	// 2. 根据key获取value
	r.key = target.URL.Path
	if err = r.sync(); err != nil {
		return nil, err
	}
	go r.watch()
	return nil, nil
}

func (r Resolver) Scheme() string {
	return "etcd"
}

func (r Resolver) sync() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(r.conf.RWTimeout)*time.Second)
	defer cancelFunc()
	// 前缀查找
	res, err := r.etcdCli.Get(ctx, r.key, clientv3.WithPrefix())
	if err != nil {
		logs.Error("grpc client get etcd failed, name=%s,err:%v", r.key, err)
		return err
	}
	r.srvAddrList = []resolver.Address{}
	for _, v := range res.Kvs {
		server, err := ParseValue(v.Value)
		if err != nil {
			logs.Error("grpc client parse etcd value failed, name=%s,err:%v", r.key, err)
			continue
		}
		r.srvAddrList = append(r.srvAddrList, resolver.Address{
			Addr:       server.Addr,
			Attributes: attributes.New("weight", server.Weight),
		})
	}
	// 告知grpc
	err = r.cc.UpdateState(resolver.State{
		Addresses: r.srvAddrList,
	})
	if err != nil {
		logs.Error("grpc client updated failed, name=%s,err:%v", r.key, err)
		return err
	}
	return nil
}

func (r Resolver) watch() {
	// 1. 定时同步数据
	// 2. 监听节点的事件，从而触发不同的操作
	// 3. 监听close事件，关闭etcd
	r.watchCh = r.etcdCli.Watch(context.Background(), r.key, clientv3.WithPrefix())
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-r.closeCh:
			r.Close()
		case res, ok := <-r.watchCh:
			if ok {
				//
				r.update(res.Events)
			}
		case <-ticker.C:
			if err := r.sync(); err != nil {
				logs.Error("watch sync failed,err :%v", err)
			}
		}
	}

}

func (r Resolver) update(events []*clientv3.Event) {
	for _, event := range events {
		switch event.Type {
		case clientv3.EventTypePut:
			server, err := ParseValue(event.Kv.Value)
			if err != nil {
				logs.Error("grpc client update(EventTypePut) etcd value failed, name=%s,err:%v", r.key, err)
			}
			addr := resolver.Address{
				Addr:       server.Addr,
				Attributes: attributes.New("weight", server.Weight),
			}
			if !Exist(r.srvAddrList, addr) {
				r.srvAddrList = append(r.srvAddrList, addr)
				err = r.cc.UpdateState(resolver.State{
					Addresses: r.srvAddrList,
				})
				if err != nil {
					logs.Error("grpc client updated(EventTypePut) failed, name=%s,err:%v", r.key, err)
				}
			}

		case clientv3.EventTypeDelete:
			// 接收到delete操作，删除r.srvAddrList
			server, err := ParseKey(string(event.Kv.Key))
			if err != nil {
				logs.Error("grpc client update(EventTypeDelete) etcd value failed, name=%s,err:%v", r.key, err)
			}
			addr := resolver.Address{
				Addr: server.Addr,
			}
			if list, ok := Remove(r.srvAddrList, addr); ok {
				r.srvAddrList = list
				err = r.cc.UpdateState(resolver.State{
					Addresses: r.srvAddrList,
				})
				if err != nil {
					logs.Error("grpc client updated(EventTypeDelete) failed, name=%s,err:%v", r.key, err)
				}
			}
		}
	}
}

func (r Resolver) Close() {
	if r.etcdCli != nil {
		err := r.etcdCli.Close()
		if err != nil {
			logs.Error("resolver close etcd error :%v", err)
		}
	}

}

func Exist(list []resolver.Address, addr resolver.Address) bool {
	for i := range list {
		if list[i].Addr == addr.Addr {
			return true
		}
	}
	return false
}

func Remove(list []resolver.Address, addr resolver.Address) ([]resolver.Address, bool) {
	for i := range list {
		if list[i].Addr == addr.Addr {
			list[i] = list[len(list)-1]
			return list[:len(list)-1], true
		}
	}
	return nil, false
}

func NewResolver(conf config.EtcdConf) *Resolver {
	return &Resolver{
		conf:        conf,
		DialTimeout: conf.DialTimeout,
	}
}
