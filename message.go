package rpcxr

// 游戏服代理请求
type RPCProxyRequest struct {
    PlayerId string
    Session  string
    Data     []byte
    Ext      map[string]interface{} `json:"ext"`
}

// 游戏服代理回复
type RPCProxyReply struct {
    Data []byte
}
