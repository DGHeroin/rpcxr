package main

import (
    "context"
    "github.com/DGHeroin/rpcxr"
    "log"
    "os"
)

func main() {
    addrs := []string{
        "127.0.0.1:20001",
        "127.0.0.1:20002",
    }
    cli := rpcxr.NewClientMultiple(addrs, "s")
    {
        var (
            r string = os.Getenv("name")
            w string
        )
        err := cli.Call(context.Background(), "mul", &r, &w)
        if err != nil {
            log.Println(err)
        }
        log.Printf("请求回复:%v", w)
    }
    {
        var (
            r string = os.Getenv("name")
            w string
        )
        err := cli.Broadcast(context.Background(), "broadcast", &r, &w)
        if err != nil {
            log.Println(err)
        }
        log.Printf("广播回复:%v", w)
    }

}
