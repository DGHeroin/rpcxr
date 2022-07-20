package rpcxr

import (
    "github.com/rpcxio/libkv/store"
    "github.com/smallnest/rpcx/client"
    "sync"
)

type (
    ClientConn struct {
        Client    client.XClient
        Discovery client.ServiceDiscovery
    }
    Client client.XClient
)

func (o *ClientConn) UpdateAddress(addr []string) bool {
    var ptr interface{} = o.Discovery
    switch u := ptr.(type) {
    case *client.MultipleServersDiscovery:
        u.Update(ParseAddress(addr))
        return true
    }
    return false
}

var (
    mu      sync.RWMutex
    clients = map[string]*ClientConn{}
)

func GetClientMultiple(servicePath string, addr ...string) *ClientConn {
    mu.RLock()
    obj, ok := clients[servicePath]
    mu.RUnlock()
    if ok {
        return obj
    }

    dis, err := client.NewMultipleServersDiscovery(ParseAddress(addr))
    if err != nil {
        return nil
    }
    p := client.NewXClient(servicePath, client.Failtry, client.ConsistentHash, dis, client.DefaultOption)
    obj = &ClientConn{
        Client:    p,
        Discovery: dis,
    }
    mu.Lock()
    clients[servicePath] = obj
    mu.Unlock()
    return obj

}

func ParseAddress(ss []string) (kvs []*client.KVPair) {
    for _, addr := range ss {
        kvs = append(kvs, &client.KVPair{Key: addr})
    }
    return
}
func NewMultipleServersDiscovery(addrs []string) client.ServiceDiscovery {
    var kvs = ParseAddress(addrs)
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

func NewClientMultiple(addr []string, ServicePath string) (client.XClient, *client.MultipleServersDiscovery) {
    discovery := NewMultipleServersDiscovery(addr)
    return client.NewXClient(ServicePath, client.Failtry, client.ConsistentHash, discovery, client.DefaultOption), discovery.(*client.MultipleServersDiscovery)
}
