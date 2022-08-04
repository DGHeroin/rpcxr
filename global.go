package rpcxr

import (
    "github.com/smallnest/rpcx/client"
    "sync"
)

var (
    mu      sync.RWMutex
    clients = map[string]*ClientConn{}
)

func GetClientMultiple(servicePath string, addr ...string) *ClientConn {
    obj, ok := getClient(servicePath)
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
    setClient(servicePath, obj)
    return obj
}
func getClient(servicePath string) (*ClientConn, bool) {
    mu.RLock()
    defer mu.RUnlock()
    obj, ok := clients[servicePath]
    return obj, ok
}
func setClient(servicePath string, client *ClientConn) {
    mu.Lock()
    defer mu.Unlock()
    clients[servicePath] = client
}
func HasClient(servicePath string) bool {
    _, ok := getClient(servicePath)
    return ok
}
