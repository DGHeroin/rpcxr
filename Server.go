package rpcxr

import (
    "github.com/smallnest/rpcx/server"
    "time"
)

func NewServer() *server.Server {
    return server.NewServer()
}

type IServer interface {
    Close()
    Server() *server.Server
}
type ServerOption struct {
    Address        string
    Servers        []string
    ServePath      string
    Zone           string
    UpdateInterval time.Duration
    Username       string
    Password       string
}
