package main

import (
    "context"
    "github.com/DGHeroin/rpcxr"
    "log"
    "os"
)

func main() {
    s := rpcxr.NewServer()
    s.RegisterFunctionName("s", "mul", onMul, "")
    s.RegisterFunctionName("s", "broadcast", onBroadcast, "")
    s.Serve("tcp", os.Getenv("addr"))

}
func onMul(ctx context.Context, r, w *string) error {
    log.Println("收到", *r)
    *w = os.Getenv("name")
    return nil
}
func onBroadcast(ctx context.Context, r, w *string) error {
    log.Println("收到广播", *r)
    *w = "广播-" + os.Getenv("name")
    return nil
}
