package rpcxr

import (
    "github.com/rpcxio/libkv/store"
    "github.com/smallnest/rpcx/client"
    "sync"
)

type ClientOption struct {
    ServicePath string
    BasePath    string
    FailType    int
    Type        string
    Address     []string
    Username    string
    Password    string
    Select      int
}

type Client client.XClient

var (
    lockClient sync.Map
)

func GetRPCClient(servicePath string, addr string) Client {
    obj, ok := lockClient.Load(addr)
    if !ok {
        dis, err := client.NewPeer2PeerDiscovery(addr, "")
        if err != nil {
            return nil
        }
        p := client.NewXClient(servicePath, client.Failfast, client.ConsistentHash, dis, client.DefaultOption)
        lockClient.Store(addr, p)
        return p
    } else {
        return obj.(Client)
    }
}

func NewMultipleServersDiscovery(addrs []string) client.ServiceDiscovery {
    var kvs []*client.KVPair
    for _, addr := range addrs {
        kvs = append(kvs, &client.KVPair{Key: addr})
    }
    dis, _ := client.NewMultipleServersDiscovery(kvs)
    return dis
}

func NewRedisDiscovery(servicePath, basePath, username, password string, servers []string) client.ServiceDiscovery {
    dis, _ := client.NewRedisDiscovery(basePath, servicePath, servers, &store.Config{
        Username: username,
        Password: password,
    })
    return dis
}
func NewEtcdv3Discovery(path string, path2 string, username string, password string, address []string) client.ServiceDiscovery {
    panic("not impl")
}
func NewXClient(opt *ClientOption) client.XClient {
    var d client.ServiceDiscovery
    switch opt.Type {
    case "multiple":
        d = NewMultipleServersDiscovery(opt.Address)
    case "redis":
        d = NewRedisDiscovery(opt.ServicePath, opt.BasePath, opt.Username, opt.Password, opt.Address)
    case "etcdv3":
        d = NewEtcdv3Discovery(opt.ServicePath, opt.BasePath, opt.Username, opt.Password, opt.Address)
    }
    cli := client.NewXClient(opt.ServicePath, client.FailMode(opt.FailType), client.SelectMode(opt.Select), d, client.DefaultOption)
    return cli
}

func NewClientMultiple(addr []string, ServicePath string) client.XClient {
    discovery := NewMultipleServersDiscovery(addr)
    return client.NewXClient(ServicePath, client.Failtry, client.ConsistentHash, discovery, client.DefaultOption)
}
