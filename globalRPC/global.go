package globalRPC

import (
    "github.com/DGHeroin/rpcxr"
    "github.com/smallnest/rpcx/client"
    "sync"
)

var (
    mu      sync.RWMutex
    clients = map[string]*rpcxr.ClientConn{}
)

type (
    option struct {
        Retries      int
        Address      []string
        FailMode     client.FailMode
        SelectMode   client.SelectMode
        customOption *client.Option
    }
    Option func(*option)
)

func GetClientMultiple(servicePath string, addr ...string) *rpcxr.ClientConn {
    obj, ok := getClient(servicePath)
    if ok {
        return obj
    }

    dis, err := client.NewMultipleServersDiscovery(rpcxr.ParseAddress(addr))
    if err != nil {
        return nil
    }
    p := client.NewXClient(servicePath, client.Failtry, client.ConsistentHash, dis, client.DefaultOption)
    obj = &rpcxr.ClientConn{
        Client:    p,
        Discovery: dis,
    }
    setClient(servicePath, obj)
    return obj
}
func getClient(servicePath string) (*rpcxr.ClientConn, bool) {
    mu.RLock()
    defer mu.RUnlock()
    obj, ok := clients[servicePath]
    return obj, ok
}
func setClient(servicePath string, client *rpcxr.ClientConn) {
    mu.Lock()
    defer mu.Unlock()
    clients[servicePath] = client
}
func HasClient(servicePath string) bool {
    _, ok := getClient(servicePath)
    return ok
}
func GetClientWithOption(servicePath string, opts ...Option) *rpcxr.ClientConn {
    obj, ok := getClient(servicePath)
    if ok {
        return obj
    }
    o := defaultOption()
    for _, opt := range opts {
        opt(o)
    }
    dis, err := client.NewMultipleServersDiscovery(rpcxr.ParseAddress(o.Address))
    if err != nil {
        return nil
    }
    dfOpt := client.DefaultOption
    dfOpt.Retries = o.Retries
    if o.customOption != nil {
        dfOpt = *o.customOption
    }
    p := client.NewXClient(servicePath, o.FailMode, o.SelectMode, dis, dfOpt)
    obj = &rpcxr.ClientConn{
        Client:    p,
        Discovery: dis,
    }
    setClient(servicePath, obj)
    return obj
}
func defaultOption() *option {
    return &option{
        Retries: 3,
    }
}
func WithAddress(addr []string) Option {
    return func(o *option) {
        o.Address = addr
    }
}
func WithRetries(retries int) Option {
    return func(o *option) {
        o.Retries = retries
    }
}
func WithCustomOption(opt *client.Option) Option {
    return func(o *option) {
        o.customOption = opt
    }
}
func WithFailMode(mode int) Option {
    return func(o *option) {
        o.FailMode = client.FailMode(mode)
    }
}
func WithSelectMode(mode int) Option {
    return func(o *option) {
        o.SelectMode = client.SelectMode(mode)
    }
}
