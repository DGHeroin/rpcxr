package main

import (
    "bufio"
    "context"
    "github.com/DGHeroin/rpcxr"
    "github.com/DGHeroin/rpcxr/global"
    "log"
    "os"
    "strings"
    "time"
)

func main() {
    cc := global.GetClientMultiple("s")
    cc.UpdateAddress([]string{"127.0.0.1:20001"})
}
func main2() {

    addrs := []string{
        // "127.0.0.1:20001",
    }
    cli, dis := rpcxr.NewClientMultiple(addrs, "s")
    go func() {
        r := bufio.NewReader(os.Stdin)

        for {
            line, _, err := r.ReadLine()
            if err == nil {
                str := strings.TrimSpace(string(line))
                if str == "" {
                    continue
                }
                infos := strings.Split(str, ",")
                log.Println("新地址:", infos)
                newAddr := rpcxr.ParseAddress(infos)

                dis.Update(newAddr)
            }
        }
    }()
    for {
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
        time.Sleep(time.Second * 2)
    }
}
